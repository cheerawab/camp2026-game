package system

import (
	"net/http"

	"github.com/sitcon-tw/camp2026-game/internal/http/httpx"
)

type PingResponse struct {
	Message      string `json:"message" example:"pong"`
	DatabaseTime string `json:"databaseTime,omitempty" example:"2026-05-29 12:00:00.000000+00"`
}

// Ping godoc
// @Summary Database ping
// @Description Runs a minimal sqlc query to confirm database access.
// @Tags system
// @Produce json
// @Success 200 {object} PingResponse
// @Failure 503 {object} httpx.ProblemDetails
// @Router /ping [get]
func (h *Handler) Ping(w http.ResponseWriter, r *http.Request) {
	if h.queries == nil {
		httpx.WriteProblem(w, r, httpx.ServiceUnavailable("database queries are not configured"))
		return
	}

	databaseTime, err := h.queries.GetDatabaseTime(r.Context())
	if err != nil {
		httpx.WriteProblem(w, r, httpx.ServiceUnavailable("database query failed"))
		return
	}

	httpx.WriteJSON(w, http.StatusOK, PingResponse{
		Message:      "pong",
		DatabaseTime: databaseTime,
	})
}
