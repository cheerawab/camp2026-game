package me

import (
	"context"
	"net/http"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/sitcon-tw/camp2026-game/internal/http/authctx"
	"github.com/sitcon-tw/camp2026-game/internal/http/httpx"
	mongomodel "github.com/sitcon-tw/camp2026-game/internal/mongodb/model"
)

// Status godoc
// @Summary Get current player status
// @Description Returns the authenticated player's profile summary, team, and open power total.
// @Tags me
// @Produce json
// @Security AuthCookieAuth
// @Success 200 {object} StatusResponse
// @Failure 401 {object} httpx.ProblemDetails
// @Failure 500 {object} httpx.ProblemDetails
// @Failure 503 {object} httpx.ProblemDetails
// @Router /me/status [get]
func (h *Handler) Status(w http.ResponseWriter, r *http.Request) {
	player, ok := currentPlayer(w, r)
	if !ok {
		return
	}
	if h.db == nil {
		httpx.WriteProblem(w, r, httpx.ServiceUnavailable("database is unavailable"))
		return
	}

	openPower, err := h.sumOpenPower(r.Context(), player.ID)
	if err != nil {
		httpx.WriteProblem(w, r, httpx.NewError(http.StatusInternalServerError, "status unavailable"))
		return
	}

	var team *mongomodel.Team
	var teamMembers []mongomodel.Player
	teamID := playerTeamID(player)
	if teamID != "" {
		foundTeam, err := h.findTeam(r.Context(), teamID)
		if err != nil {
			httpx.WriteProblem(w, r, httpx.NewError(http.StatusInternalServerError, "status unavailable"))
			return
		}
		team = &foundTeam

		teamMembers, err = h.findTeamMembers(r.Context(), teamID)
		if err != nil {
			httpx.WriteProblem(w, r, httpx.NewError(http.StatusInternalServerError, "status unavailable"))
			return
		}
	}

	httpx.WriteJSON(w, http.StatusOK, statusResponse(player, team, openPower, teamMembers))
}

func playerTeamID(player mongomodel.Player) string {
	if player.Role == authctx.PlayerRoleStaff {
		return ""
	}
	return player.TeamID
}

func (h *Handler) findTeam(ctx context.Context, teamID string) (mongomodel.Team, error) {
	var team mongomodel.Team
	err := h.db.Collection(mongomodel.TeamsCollection).
		FindOne(ctx, bson.M{"_id": teamID}).
		Decode(&team)
	if err != nil {
		return mongomodel.Team{}, err
	}
	if team.ID == "" || team.Name == "" {
		return mongomodel.Team{}, mongo.ErrNoDocuments
	}
	return team, nil
}

func (h *Handler) findTeamMembers(ctx context.Context, teamID string) ([]mongomodel.Player, error) {
	cursor, err := h.db.Collection(mongomodel.PlayersCollection).Find(
		ctx,
		bson.M{
			"team_id": teamID,
			"role":    bson.M{"$ne": authctx.PlayerRoleStaff},
		},
		options.Find().
			SetProjection(bson.D{
				{Key: "auth_token", Value: 0},
				{Key: "qrcode_token", Value: 0},
				{Key: "default_sitone_ids", Value: 0},
			}).
			SetSort(bson.D{
				{Key: "nickname", Value: 1},
				{Key: "_id", Value: 1},
			}),
	)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = cursor.Close(ctx)
	}()

	var players []mongomodel.Player
	if err := cursor.All(ctx, &players); err != nil {
		return nil, err
	}
	if players == nil {
		return []mongomodel.Player{}, nil
	}
	return players, nil
}

func (h *Handler) sumOpenPower(ctx context.Context, playerID string) (int, error) {
	cursor, err := h.db.Collection(mongomodel.OpenPowerRecordsCollection).
		Aggregate(ctx, openPowerTotalPipeline(playerID))
	if err != nil {
		return 0, err
	}
	defer func() {
		_ = cursor.Close(ctx)
	}()

	return openPowerTotalFromCursor(ctx, cursor)
}

func openPowerTotalFromCursor(ctx context.Context, cursor *mongo.Cursor) (int, error) {
	var totals []struct {
		Total int `bson:"total"`
	}
	if err := cursor.All(ctx, &totals); err != nil {
		return 0, err
	}
	if len(totals) == 0 {
		return 0, nil
	}
	return totals[0].Total, nil
}

func openPowerTotalPipeline(playerID string) mongo.Pipeline {
	return mongo.Pipeline{
		{{Key: "$match", Value: bson.D{{Key: "player_id", Value: playerID}}}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: nil},
			{Key: "total", Value: bson.D{{Key: "$sum", Value: "$amount"}}},
		}}},
	}
}

func statusResponse(player mongomodel.Player, team *mongomodel.Team, openPower int, teamMembers []mongomodel.Player) StatusResponse {
	response := StatusResponse{
		PlayerID:    player.ID,
		Nickname:    player.Nickname,
		TeamMembers: teamMemberResponses(teamMembers),
		OpenPower:   openPower,
		AvatarURL:   player.AvatarURL,
		Role:        player.Role,
	}
	if team != nil {
		response.Team = &TeamResponse{
			TeamID: team.ID,
			Name:   team.Name,
		}
	}
	return response
}

func teamMemberResponses(players []mongomodel.Player) []TeamMemberResponse {
	members := make([]TeamMemberResponse, 0, len(players))
	for _, player := range players {
		if player.ID == "" || player.Nickname == "" || player.Role == authctx.PlayerRoleStaff {
			continue
		}
		members = append(members, TeamMemberResponse{
			PlayerID:  player.ID,
			Nickname:  player.Nickname,
			AvatarURL: player.AvatarURL,
			Role:      player.Role,
		})
	}
	return members
}
