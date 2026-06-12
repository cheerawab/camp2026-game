package matches

import (
	"errors"
	"net/http"
	"strings"

	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/sitcon-tw/camp2026-game/internal/http/httpx"
	mongomodel "github.com/sitcon-tw/camp2026-game/internal/mongodb/model"
)

// Join godoc
// @Summary Join match room
// @Description Joins a waiting two-player quiz match room by invite code.
// @Tags matches
// @Accept json
// @Produce json
// @Security AuthCookieAuth
// @Param request body JoinMatchRequest true "Join match request"
// @Success 200 {object} JoinMatchResponse
// @Failure 400 {object} httpx.ProblemDetails
// @Failure 401 {object} httpx.ProblemDetails
// @Failure 404 {object} httpx.ProblemDetails
// @Failure 409 {object} httpx.ProblemDetails
// @Failure 422 {object} httpx.ProblemDetails
// @Failure 500 {object} httpx.ProblemDetails
// @Failure 503 {object} httpx.ProblemDetails
// @Router /matches/join [post]
func (h *Handler) Join(w http.ResponseWriter, r *http.Request) {
	player, ok := currentPlayer(w, r)
	if !ok || !h.requireDatabase(w, r) || !h.requireContent(w, r) {
		return
	}

	var body JoinMatchRequest
	if err := httpx.DecodeJSON(r, &body); err != nil {
		httpx.WriteProblem(w, r, err)
		return
	}
	body.Code = strings.ToUpper(strings.TrimSpace(body.Code))
	if err := httpx.ValidateStruct(body); err != nil {
		httpx.WriteProblem(w, r, err)
		return
	}

	match, err := h.findMatchByCode(r.Context(), body.Code)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			httpx.WriteProblem(w, r, httpx.NotFound("match not found"))
			return
		}
		httpx.WriteProblem(w, r, httpx.InternalServerError("match join failed", "match_join_lookup_failed", err))
		return
	}

	h.joinMatch(w, r, match, player)
}

func (h *Handler) joinMatch(w http.ResponseWriter, r *http.Request, match mongomodel.Match, player mongomodel.Player) {
	if match.Status != mongomodel.MatchStatusWaiting {
		httpx.WriteProblem(w, r, httpx.NewError(http.StatusConflict, "match is not joinable"))
		return
	}
	if isParticipant(match, player.ID) {
		httpx.WriteProblem(w, r, httpx.NewError(http.StatusConflict, "player already joined match"))
		return
	}
	if len(match.Players) >= 2 {
		httpx.WriteProblem(w, r, httpx.NewError(http.StatusConflict, "match is full"))
		return
	}
	sitoneIDs, err := h.defaultSitoneLoadout(r.Context(), player)
	if err != nil {
		httpx.WriteProblem(w, r, httpx.InternalServerError("match join failed", "match_join_default_loadout_failed", err))
		return
	}

	match.Players = append(match.Players, mongomodel.MatchPlayer{
		PlayerID:  player.ID,
		Nickname:  player.Nickname,
		Ready:     false,
		Score:     0,
		SitoneIDs: sitoneIDs,
	})
	if err := h.saveMatch(r.Context(), match); err != nil {
		httpx.WriteProblem(w, r, httpx.InternalServerError("match join failed", "match_join_save_failed", err))
		return
	}

	h.publishState(r.Context(), match, "match_updated")
	state, err := h.buildMatchState(r.Context(), match)
	if err != nil {
		httpx.WriteProblem(w, r, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, state)
}
