package matches

import (
	"errors"
	"math/rand"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/sitcon-tw/camp2026-game/internal/http/httpx"
	mongomodel "github.com/sitcon-tw/camp2026-game/internal/mongodb/model"
)

// Ready godoc
// @Summary Mark current player ready
// @Description Marks the authenticated player ready. The match starts automatically when both players are ready.
// @Tags matches
// @Produce json
// @Security AuthCookieAuth
// @Success 200 {object} ReadyMatchResponse
// @Failure 401 {object} httpx.ProblemDetails
// @Failure 404 {object} httpx.ProblemDetails
// @Failure 409 {object} httpx.ProblemDetails
// @Failure 500 {object} httpx.ProblemDetails
// @Failure 503 {object} httpx.ProblemDetails
// @Router /matches/{matchID}/ready [post]
func (h *Handler) Ready(w http.ResponseWriter, r *http.Request) {
	player, ok := currentPlayer(w, r)
	if !ok || !h.requireDatabase(w, r) || !h.requireContent(w, r) {
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
		httpx.WriteProblem(w, r, httpx.NewError(http.StatusConflict, "match is not waiting for ready"))
		return
	}

	idx := playerIndex(match, player.ID)
	if len(match.Players[idx].SitoneIDs) == 0 {
		sitoneIDs, err := h.defaultSitoneLoadout(r.Context(), player)
		if err != nil {
			httpx.WriteProblem(w, r, httpx.InternalServerError("ready failed", "match_ready_default_loadout_failed", err))
			return
		}
		match.Players[idx].SitoneIDs = sitoneIDs
	}
	if len(match.Players[idx].SitoneIDs) == 0 {
		httpx.WriteProblem(w, r, httpx.NewError(http.StatusConflict, "select at least one sitone before ready"))
		return
	}
	match.Players[idx].Ready = true
	events := []string{"player_ready"}

	if allPlayersReady(match) {
		questionIDs, err := h.pickQuestionIDs()
		if err != nil {
			httpx.WriteProblem(w, r, httpx.InternalServerError("match start failed", "match_ready_pick_questions_failed", err))
			return
		}
		now := time.Now()
		match.Status = mongomodel.MatchStatusActive
		match.Phase = mongomodel.MatchPhaseAnswering
		match.QuestionIDs = questionIDs
		match.CurrentQuestionIndex = 0
		match.StartedAt = now
		match.RoundStartedAt = now
		match.RoundEndsAt = now.Add(roundDuration * time.Second)
		events = append(events, "round_started")
	}

	if err := h.saveMatch(r.Context(), match); err != nil {
		httpx.WriteProblem(w, r, httpx.InternalServerError("ready failed", "match_ready_save_failed", err))
		return
	}

	for _, event := range events {
		h.publishState(r.Context(), match, event)
	}
	state, err := h.buildMatchState(r.Context(), match)
	if err != nil {
		httpx.WriteProblem(w, r, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, state)
}

func (h *Handler) pickQuestionIDs() ([]string, error) {
	questions := h.content.ListQuizQuestions()
	if len(questions) < matchQuestionCount {
		return nil, httpx.InternalServerError("quiz questions are unavailable", "match_questions_insufficient", errors.New("not enough quiz questions"))
	}

	rand.Shuffle(len(questions), func(i, j int) {
		questions[i], questions[j] = questions[j], questions[i]
	})

	ids := make([]string, 0, matchQuestionCount)
	for _, question := range questions[:matchQuestionCount] {
		ids = append(ids, question.ID)
	}
	return ids, nil
}
