package worldboss

import (
	"net/http"

	"github.com/sitcon-tw/camp2026-game/internal/http/apimodel"
	"github.com/sitcon-tw/camp2026-game/internal/http/httpx"
)

// GetBoss godoc
// @Summary Get world boss
// @Description Returns a single world boss, current stage progress, remaining attempts, and reward state.
// @Tags WorldBoss
// @Produce json
// @Security AuthCookieAuth
// @Param bossID path string true "World boss ID"
// @Success 200 {object} apimodel.WorldBossDetailResponse
// @Failure 401 {object} httpx.ProblemDetails
// @Failure 404 {object} httpx.ProblemDetails
// @Failure 500 {object} httpx.ProblemDetails
// @Router /world-bosses/{bossID} [get]
func (h *Handler) GetBoss(w http.ResponseWriter, _ *http.Request) {
	httpx.WriteJSON(w, http.StatusOK, apimodel.WorldBossDetailResponse{
		Boss: apimodel.WorldBossSummary{
			BossID:            "boss_layer_1",
			Name:              "Knowledge Core",
			Layer:             1,
			HP:                4500,
			MaxHP:             10000,
			RemainingAttempts: 2,
			Status:            "active",
		},
		Rewards: []apimodel.WorldBossReward{
			{
				RewardID: "wb_reward_layer_1",
				Status:   "locked",
				Reward:   apimodel.Reward{OpenPower: 200},
			},
		},
	})
}
