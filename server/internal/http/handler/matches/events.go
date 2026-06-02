package matches

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/sitcon-tw/camp2026-game/internal/http/httpx"
)

// Events godoc
// @Summary Stream match events
// @Description Streams match state changes for participants with Server-Sent Events.
// @Tags matches
// @Produce text/event-stream
// @Security AuthCookieAuth
// @Success 200 {string} string "SSE event stream"
// @Failure 401 {object} httpx.ProblemDetails
// @Failure 404 {object} httpx.ProblemDetails
// @Failure 500 {object} httpx.ProblemDetails
// @Failure 503 {object} httpx.ProblemDetails
// @Router /matches/{matchID}/events [get]
func (h *Handler) Events(w http.ResponseWriter, r *http.Request) {
	player, ok := currentPlayer(w, r)
	if !ok || !h.requireDatabase(w, r) || !h.requireContent(w, r) {
		return
	}

	matchID := chi.URLParam(r, "matchID")
	match, err := h.findMatchByID(r.Context(), matchID)
	if err != nil {
		writeMatchProblem(w, r, err)
		return
	}
	if !isParticipant(match, player.ID) {
		httpx.WriteProblem(w, r, httpx.NotFound("match not found"))
		return
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		httpx.WriteProblem(w, r, httpx.NewError(http.StatusInternalServerError, "event stream is unavailable"))
		return
	}

	events, unsubscribe := h.broker.Subscribe(matchID)
	defer unsubscribe()

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	state, err := h.buildMatchState(r.Context(), match)
	if err != nil {
		httpx.WriteProblem(w, r, err)
		return
	}
	writeSSE(w, "match_updated", state)
	flusher.Flush()

	heartbeat := time.NewTicker(15 * time.Second)
	defer heartbeat.Stop()

	for {
		select {
		case <-r.Context().Done():
			return
		case event, ok := <-events:
			if !ok {
				return
			}
			writeSSE(w, event.Name, event.Data)
			flusher.Flush()
		case <-heartbeat.C:
			_, _ = fmt.Fprint(w, ": keepalive\n\n")
			flusher.Flush()
		}
	}
}

func writeSSE(w http.ResponseWriter, name string, data MatchStateResponse) {
	payload, err := json.Marshal(data)
	if err != nil {
		return
	}
	_, _ = fmt.Fprintf(w, "event: %s\n", name)
	_, _ = fmt.Fprintf(w, "data: %s\n\n", payload)
}
