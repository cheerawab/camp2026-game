package me

import (
	"context"
	"net/http"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

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

	team, err := h.findTeam(r.Context(), player.TeamID)
	if err != nil {
		httpx.WriteProblem(w, r, httpx.NewError(http.StatusInternalServerError, "status unavailable"))
		return
	}

	openPower, err := h.sumOpenPower(r.Context(), player.ID)
	if err != nil {
		httpx.WriteProblem(w, r, httpx.NewError(http.StatusInternalServerError, "status unavailable"))
		return
	}

	httpx.WriteJSON(w, http.StatusOK, statusResponse(player, team, openPower))
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

func statusResponse(player mongomodel.Player, team mongomodel.Team, openPower int) StatusResponse {
	return StatusResponse{
		PlayerID: player.ID,
		Nickname: player.Nickname,
		Team: TeamResponse{
			TeamID: team.ID,
			Name:   team.Name,
		},
		OpenPower: openPower,
		AvatarURL: player.AvatarURL,
		Role:      player.Role,
	}
}
