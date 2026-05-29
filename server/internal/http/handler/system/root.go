package system

import (
	"net/http"

	"github.com/sitcon-tw/camp2026-game/internal/http/httpx"
)

type MessageResponse struct {
	Message string `json:"message" example:"Camp 2026 Game API is running"`
}

// Root godoc
// @Summary Hello world
// @Description Confirms the API process is running.
// @Tags system
// @Produce json
// @Success 200 {object} MessageResponse
// @Router / [get]
func (h *Handler) Root(w http.ResponseWriter, _ *http.Request) {
	httpx.WriteJSON(w, http.StatusOK, MessageResponse{
		Message: "Camp 2026 Game API is running",
	})
}
