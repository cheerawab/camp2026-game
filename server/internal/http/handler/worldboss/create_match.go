package worldboss

import (
	"net/http"

	"github.com/sitcon-tw/camp2026-game/internal/http/apimodel"
	"github.com/sitcon-tw/camp2026-game/internal/http/httpx"
)

// CreateMatch godoc
// @Summary Create world boss match
// @Description Starts a world boss match and consumes one player attempt if accepted.
// @Tags WorldBoss
// @Accept json
// @Produce json
// @Security AuthCookieAuth
// @Param bossID path string true "World boss ID"
// @Param request body apimodel.WorldBossMatchCreateRequest true "World boss match request"
// @Success 201 {object} apimodel.WorldBossMatchResponse
// @Failure 400 {object} httpx.ProblemDetails
// @Failure 401 {object} httpx.ProblemDetails
// @Failure 422 {object} httpx.ProblemDetails
// @Failure 501 {object} httpx.ProblemDetails
// @Router /world-bosses/{bossID}/matches [post]
func (h *Handler) CreateMatch(w http.ResponseWriter, r *http.Request) {
	var body apimodel.WorldBossMatchCreateRequest
	if err := httpx.DecodeJSON(r, &body); err != nil {
		httpx.WriteProblem(w, r, err)
		return
	}
	if err := httpx.ValidateStruct(body); err != nil {
		httpx.WriteProblem(w, r, err)
		return
	}

	httpx.WriteProblem(w, r, httpx.NotImplemented("world boss match creation is not implemented yet"))
}
