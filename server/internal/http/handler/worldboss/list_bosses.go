package worldboss

import (
	"net/http"

	"github.com/sitcon-tw/camp2026-game/internal/http/apimodel"
	"github.com/sitcon-tw/camp2026-game/internal/http/httpx"
)

// ListBosses godoc
// @Summary List world bosses
// @Description Lists active world bosses, shared progress, and the current player's remaining challenge attempts.
// @Tags WorldBoss
// @Produce json
// @Security AuthCookieAuth
// @Success 200 {object} apimodel.WorldBossListResponse
// @Failure 401 {object} httpx.ProblemDetails
// @Failure 500 {object} httpx.ProblemDetails
// @Router /world-bosses [get]
func (h *Handler) ListBosses(w http.ResponseWriter, _ *http.Request) {
	httpx.WriteJSON(w, http.StatusOK, apimodel.WorldBossListResponse{
		Bosses: []apimodel.WorldBossSummary{
			{
				BossID:            "boss_layer_1",
				Name:              "Knowledge Core",
				Layer:             1,
				HP:                4500,
				MaxHP:             10000,
				RemainingAttempts: 2,
				Status:            "active",
			},
		},
	})
}
