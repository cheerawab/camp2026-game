package qr

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/sitcon-tw/camp2026-game/internal/http/httpx"
	mongomodel "github.com/sitcon-tw/camp2026-game/internal/mongodb/model"
)

// Resolve godoc
// @Summary Resolve player QR code
// @Description Resolves a player QR code identifier into a public player summary without exposing auth credentials.
// @Tags qr
// @Accept json
// @Produce json
// @Param request body ResolveRequest true "QR resolve request"
// @Success 200 {object} ResolveResponse
// @Failure 400 {object} httpx.ProblemDetails
// @Failure 404 {object} httpx.ProblemDetails
// @Failure 422 {object} httpx.ProblemDetails
// @Failure 500 {object} httpx.ProblemDetails
// @Failure 503 {object} httpx.ProblemDetails
// @Router /qr/resolve [post]
func (h *Handler) Resolve(w http.ResponseWriter, r *http.Request) {
	if h.db == nil {
		httpx.WriteProblem(w, r, httpx.ServiceUnavailable("database is unavailable"))
		return
	}

	var body ResolveRequest
	if err := httpx.DecodeJSON(r, &body); err != nil {
		httpx.WriteProblem(w, r, err)
		return
	}
	body.QRCodeToken = strings.TrimSpace(body.QRCodeToken)
	if err := httpx.ValidateStruct(body); err != nil {
		httpx.WriteProblem(w, r, err)
		return
	}

	player, err := h.findPlayerByQRCodeToken(r.Context(), body.QRCodeToken)
	if errors.Is(err, mongo.ErrNoDocuments) {
		httpx.WriteProblem(w, r, httpx.NotFound("qr code not found"))
		return
	}
	if err != nil {
		httpx.WriteProblem(w, r, httpx.InternalServerError("qr resolve failed", "qr_player_lookup_failed", err))
		return
	}

	team, err := h.findTeam(r.Context(), player.TeamID)
	if err != nil {
		httpx.WriteProblem(w, r, httpx.InternalServerError("qr resolve failed", "qr_team_lookup_failed", err))
		return
	}
	openPower, err := h.sumOpenPower(r.Context(), player.ID)
	if err != nil {
		httpx.WriteProblem(w, r, httpx.InternalServerError("qr resolve failed", "qr_open_power_sum_failed", err))
		return
	}

	httpx.WriteJSON(w, http.StatusOK, ResolveResponse{
		Player: PlayerSummaryResponse{
			PlayerID:  player.ID,
			Nickname:  player.Nickname,
			AvatarURL: player.AvatarURL,
			OpenPower: openPower,
			Team: TeamResponse{
				TeamID: team.ID,
				Name:   team.Name,
			},
		},
	})
}

func (h *Handler) findPlayerByQRCodeToken(ctx context.Context, token string) (mongomodel.Player, error) {
	var player mongomodel.Player
	err := h.db.Collection(mongomodel.PlayersCollection).
		FindOne(ctx, bson.M{"qrcode_token": token}).
		Decode(&player)
	return player, err
}

func (h *Handler) findTeam(ctx context.Context, teamID string) (mongomodel.Team, error) {
	var team mongomodel.Team
	err := h.db.Collection(mongomodel.TeamsCollection).
		FindOne(ctx, bson.M{"_id": teamID}).
		Decode(&team)
	return team, err
}

func (h *Handler) sumOpenPower(ctx context.Context, playerID string) (int, error) {
	cursor, err := h.db.Collection(mongomodel.OpenPowerRecordsCollection).Aggregate(ctx, mongo.Pipeline{
		{{Key: "$match", Value: bson.D{{Key: "player_id", Value: playerID}}}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: nil},
			{Key: "total", Value: bson.D{{Key: "$sum", Value: "$amount"}}},
		}}},
	})
	if err != nil {
		return 0, err
	}
	defer func() {
		_ = cursor.Close(ctx)
	}()

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
