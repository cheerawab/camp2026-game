package matches

import (
	"errors"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/sitcon-tw/camp2026-game/internal/http/httpx"
	mongomodel "github.com/sitcon-tw/camp2026-game/internal/mongodb/model"
)

// Open godoc
// @Summary Get current open match
// @Description Returns the current waiting or active match for the authenticated participant.
// @Tags matches
// @Produce json
// @Security AuthCookieAuth
// @Success 200 {object} OpenMatchResponse
// @Failure 401 {object} httpx.ProblemDetails
// @Failure 404 {object} httpx.ProblemDetails
// @Failure 500 {object} httpx.ProblemDetails
// @Failure 503 {object} httpx.ProblemDetails
// @Router /matches/open [get]
func (h *Handler) Open(w http.ResponseWriter, r *http.Request) {
	player, ok := currentPlayer(w, r)
	if !ok || !h.requireDatabase(w, r) || !h.requireContent(w, r) {
		return
	}

	match, err := h.findOpenParticipantMatch(r.Context(), player.ID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			httpx.WriteProblem(w, r, httpx.NotFound("match not found"))
			return
		}
		httpx.WriteProblem(w, r, httpx.InternalServerError("match lookup failed", "match_open_lookup_failed", err))
		return
	}

	match, events, err := h.advanceMatch(r.Context(), match, time.Now())
	if err != nil {
		httpx.WriteProblem(w, r, httpx.InternalServerError("match lookup failed", "match_open_advance_failed", err))
		return
	}
	for _, event := range events {
		h.publishState(r.Context(), match, event)
	}
	if match.Status == mongomodel.MatchStatusCompleted {
		httpx.WriteProblem(w, r, httpx.NotFound("match not found"))
		return
	}

	state, err := h.buildMatchState(r.Context(), match, player.ID)
	if err != nil {
		httpx.WriteProblem(w, r, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, state)
}
