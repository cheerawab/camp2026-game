package matches

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

// JoinByQR godoc
// @Summary Join match room by QR code
// @Description Joins a waiting match by scanning either a match invite code token or the host player's QR token.
// @Tags matches
// @Accept json
// @Produce json
// @Security AuthCookieAuth
// @Param request body JoinByQRRequest true "Join by QR request"
// @Success 200 {object} JoinMatchResponse
// @Failure 400 {object} httpx.ProblemDetails
// @Failure 401 {object} httpx.ProblemDetails
// @Failure 404 {object} httpx.ProblemDetails
// @Failure 409 {object} httpx.ProblemDetails
// @Failure 422 {object} httpx.ProblemDetails
// @Failure 500 {object} httpx.ProblemDetails
// @Failure 503 {object} httpx.ProblemDetails
// @Router /matches/join-by-qr [post]
func (h *Handler) JoinByQR(w http.ResponseWriter, r *http.Request) {
	player, ok := currentPlayer(w, r)
	if !ok || !h.requireDatabase(w, r) {
		return
	}

	var body JoinByQRRequest
	if err := httpx.DecodeJSON(r, &body); err != nil {
		httpx.WriteProblem(w, r, err)
		return
	}
	body.QRCodeToken = strings.TrimSpace(body.QRCodeToken)
	if err := httpx.ValidateStruct(body); err != nil {
		httpx.WriteProblem(w, r, err)
		return
	}

	match, err := h.findMatchByQRToken(r.Context(), body.QRCodeToken)
	if errors.Is(err, mongo.ErrNoDocuments) {
		httpx.WriteProblem(w, r, httpx.NotFound("match not found"))
		return
	}
	if err != nil {
		httpx.WriteProblem(w, r, httpx.NewError(http.StatusInternalServerError, "match join failed"))
		return
	}

	h.joinMatch(w, r, match, player)
}

func (h *Handler) findMatchByQRToken(ctx context.Context, token string) (mongomodel.Match, error) {
	if match, err := h.findMatchByCode(ctx, strings.ToUpper(token)); err == nil {
		return match, nil
	} else if !errors.Is(err, mongo.ErrNoDocuments) {
		return mongomodel.Match{}, err
	}

	var host mongomodel.Player
	err := h.db.Collection(mongomodel.PlayersCollection).
		FindOne(ctx, bson.M{"qrcode_token": token}).
		Decode(&host)
	if err != nil {
		return mongomodel.Match{}, err
	}

	var match mongomodel.Match
	err = h.db.Collection(mongomodel.MatchesCollection).
		FindOne(ctx, bson.M{
			"host_player_id": host.ID,
			"status":         mongomodel.MatchStatusWaiting,
		}).
		Decode(&match)
	return match, err
}
