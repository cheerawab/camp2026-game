package bingo

import (
	"net/http"

	"github.com/sitcon-tw/camp2026-game/internal/http/apimodel"
	"github.com/sitcon-tw/camp2026-game/internal/http/httpx"
)

// CompleteMission godoc
// @Summary Complete a bingo mission
// @Description Attempts to complete a mission by player-submitted flag or server-side condition detection.
// @Tags Bingo
// @Accept json
// @Produce json
// @Security AuthCookieAuth
// @Param missionID path string true "Mission ID"
// @Param request body apimodel.MissionCompleteRequest true "Mission completion request"
// @Success 200 {object} apimodel.MissionCompleteResponse
// @Failure 400 {object} httpx.ProblemDetails
// @Failure 401 {object} httpx.ProblemDetails
// @Failure 422 {object} httpx.ProblemDetails
// @Failure 501 {object} httpx.ProblemDetails
// @Router /bingo/missions/{missionID}/complete [post]
func (h *Handler) CompleteMission(w http.ResponseWriter, r *http.Request) {
	var body apimodel.MissionCompleteRequest
	if err := httpx.DecodeJSON(r, &body); err != nil {
		httpx.WriteProblem(w, r, err)
		return
	}
	if err := httpx.ValidateStruct(body); err != nil {
		httpx.WriteProblem(w, r, err)
		return
	}

	httpx.WriteProblem(w, r, httpx.NotImplemented("bingo mission completion is not implemented yet"))
}
