package matches

import (
	"net/http"

	"github.com/sitcon-tw/camp2026-game/internal/http/httpx"
)

// GetMatchWebSocket godoc
// @Summary Match WebSocket endpoint
// @Description Upgrades to the shared real-time match WebSocket. The event contract is shared by QR duel, offline duel, and world boss matches.
// @Tags Matches
// @Produce json
// @Security AuthCookieAuth
// @Param matchID path string true "Match ID"
// @Success 200 {object} apimodel.MatchWebSocketInfoResponse
// @Failure 401 {object} httpx.ProblemDetails
// @Failure 404 {object} httpx.ProblemDetails
// @Failure 501 {object} httpx.ProblemDetails
// @Router /matches/{matchID}/ws [get]
func (h *Handler) GetMatchWebSocket(w http.ResponseWriter, r *http.Request) {
	httpx.WriteProblem(w, r, httpx.NotImplemented("match websocket is not implemented yet"))
}
