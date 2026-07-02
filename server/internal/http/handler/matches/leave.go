package matches

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/sitcon-tw/camp2026-game/internal/http/httpx"
	mongomodel "github.com/sitcon-tw/camp2026-game/internal/mongodb/model"
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

	match, err := h.findMatchByID(r.Context(), chi.URLParam(r, "matchID"))
	if err != nil {
		writeMatchProblem(w, r, err)
		return
	}
	if !isParticipant(match, player.ID) {
		httpx.WriteProblem(w, r, httpx.NotFound("match not found"))
		return
	}
	if match.Status != mongomodel.MatchStatusWaiting {
		httpx.WriteProblem(w, r, httpx.NewError(http.StatusConflict, "match has already started"))
		return
	}

	if match.HostPlayerID == player.ID {
		result, err := h.db.Collection(mongomodel.MatchesCollection).DeleteOne(
			r.Context(),
			bson.M{
				"_id":            match.ID,
				"status":         mongomodel.MatchStatusWaiting,
				"host_player_id": player.ID,
			},
		)
		if err != nil {
			httpx.WriteProblem(w, r, httpx.InternalServerError("leave match failed", "match_leave_delete_failed", err))
			return
		}
		if result.DeletedCount == 0 {
			writeMatchProblem(w, r, errMatchSaveConflict)
			return
		}
		h.publishState(r.Context(), match, "match_deleted")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	idx := playerIndex(match, player.ID)
	if idx < 0 {
		httpx.WriteProblem(w, r, httpx.NotFound("match not found"))
		return
	}
	match.Players = append(match.Players[:idx], match.Players[idx+1:]...)
	if err := h.saveMatch(r.Context(), &match); err != nil {
		if errors.Is(err, errOpenParticipantMatchExists) || errors.Is(err, errMatchSaveConflict) || mongo.IsDuplicateKeyError(err) {
			writeMatchProblem(w, r, err)
			return
		}
		httpx.WriteProblem(w, r, httpx.InternalServerError("leave match failed", "match_leave_save_failed", err))
		return
	}
	h.publishState(r.Context(), match, "match_updated")
	w.WriteHeader(http.StatusNoContent)
}
