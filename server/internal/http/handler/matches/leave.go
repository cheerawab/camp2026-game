package matches

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/sitcon-tw/camp2026-game/internal/http/httpx"
)

// Leave godoc
// @Summary Leave waiting match room
// @Description Leaves a waiting match room. If the host leaves before the match starts, the room is deleted; otherwise the player is removed from the waiting room.
// @Tags matches
// @Produce json
// @Security AuthCookieAuth
// @Success 204
// @Failure 401 {object} httpx.ProblemDetails
// @Failure 404 {object} httpx.ProblemDetails
// @Failure 409 {object} httpx.ProblemDetails
// @Failure 500 {object} httpx.ProblemDetails
// @Failure 503 {object} httpx.ProblemDetails
// @Router /matches/{matchID}/leave [post]
func (h *Handler) Leave(w http.ResponseWriter, r *http.Request) {
	player, ok := currentPlayer(w, r)
	if !ok || !h.requireDatabase(w, r) {
		return
	}

	session, err := h.sessions.GetOrLoad(r.Context(), chi.URLParam(r, "matchID"))
	if err != nil {
		writeMatchProblem(w, r, err)
		return
	}
	if _, err := session.Leave(r.Context(), player.ID); err != nil {
		if errors.Is(err, errOpenParticipantMatchExists) || errors.Is(err, errMatchSaveConflict) || mongo.IsDuplicateKeyError(err) {
			writeMatchProblem(w, r, err)
			return
		}
		httpx.WriteProblem(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
