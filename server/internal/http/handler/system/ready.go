package system

import (
	"context"
	"net/http"
	"time"

	"github.com/sitcon-tw/camp2026-game/internal/http/httpx"
)

// Ready godoc
// @Summary Readiness check
// @Description Confirms required dependencies are reachable.
// @Tags system
// @Produce json
// @Success 200 {object} StatusResponse
// @Failure 503 {object} httpx.ProblemDetails
// @Router /readyz [get]
func (h *Handler) Ready(w http.ResponseWriter, r *http.Request) {
	if h.readinessCheck == nil {
		httpx.WriteJSON(w, http.StatusOK, StatusResponse{Status: "ok"})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	if err := h.readinessCheck(ctx); err != nil {
		httpx.WriteProblem(w, r, httpx.ServiceUnavailable("readiness check failed"))
		return
	}

	httpx.WriteJSON(w, http.StatusOK, StatusResponse{Status: "ok"})
}
