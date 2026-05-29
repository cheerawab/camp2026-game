package system

import (
	"net/http"

	"github.com/sitcon-tw/camp2026-game/internal/http/httpx"
)

type StatusResponse struct {
	Status string `json:"status" example:"ok"`
}

// Health godoc
// @Summary Health check
// @Description Confirms the HTTP process is alive.
// @Tags system
// @Produce json
// @Success 200 {object} StatusResponse
// @Router /healthz [get]
func (h *Handler) Health(w http.ResponseWriter, _ *http.Request) {
	httpx.WriteJSON(w, http.StatusOK, StatusResponse{Status: "ok"})
}
