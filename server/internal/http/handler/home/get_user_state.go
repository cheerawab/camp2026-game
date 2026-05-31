package home

import (
	"net/http"

	"github.com/sitcon-tw/camp2026-game/internal/http/apimodel"
	"github.com/sitcon-tw/camp2026-game/internal/http/httpx"
)

// GetUserState godoc
// @Summary Get current user state
// @Description Returns the current authenticated user's aggregate state for the homepage and feature pages. This response uses explicit counters instead of frontend route paths or presentation metadata.
// @Tags Users
// @Produce json
// @Security AuthCookieAuth
// @Success 200 {object} apimodel.UserStateResponse
// @Failure 401 {object} httpx.ProblemDetails
// @Failure 500 {object} httpx.ProblemDetails
// @Router /users/state [get]
func (h *Handler) GetUserState(w http.ResponseWriter, _ *http.Request) {
	httpx.WriteJSON(w, http.StatusOK, apimodel.UserStateResponse{
		Player: apimodel.UserStatePlayer{
			PlayerID:  "player_01HR9Z7E2Z2VJ2QZ4P4Z",
			TeamID:    "team_blue",
			Nickname:  "Alice",
			OpenPower: 1280,
		},
		Stats: apimodel.UserStateStats{
			OpenPower:                   1280,
			CompletedMissionCount:       7,
			ClaimableMissionRewardCount: 2,
			MatchCount:                  12,
			MatchWinCount:               8,
			OwnedSitoneCount:            5,
			OwnedItemCount:              3,
		},
		Features: apimodel.UserFeatureStates{
			Bingo: apimodel.UserBingoState{
				Enabled:               true,
				ActiveMissionCount:    4,
				CompletedMissionCount: 7,
				ClaimableRewardCount:  2,
			},
			Matches: apimodel.UserMatchState{
				Enabled:             true,
				PendingPairingCount: 1,
				ActiveMatchCount:    0,
				CompletedMatchCount: 12,
				WinCount:            8,
			},
			WorldBoss: apimodel.UserWorldBossState{
				Enabled:               true,
				ActiveBossID:          "boss_layer_1",
				RemainingAttemptCount: 2,
				ClaimableRewardCount:  0,
			},
			Storage: apimodel.UserStorageState{
				Enabled:              true,
				SitoneCount:          5,
				ItemCount:            3,
				CraftableRecipeCount: 1,
			},
			QRCode: apimodel.UserQRCodeState{
				Enabled:        true,
				HasActiveToken: true,
			},
		},
		Loadout: []apimodel.SitoneSummary{
			{
				ID:           "sitone_01HR9Z7E2Z2VJ2QZ4P4Z",
				DefinitionID: "sitone-engineering",
				Name:         "Engineering Sitone",
				Type:         "engineering",
				Rarity:       "rare",
				Style:        "default",
			},
		},
	})
}
