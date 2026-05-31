package bingo

import (
	"net/http"

	"github.com/sitcon-tw/camp2026-game/internal/http/apimodel"
	"github.com/sitcon-tw/camp2026-game/internal/http/httpx"
)

// ListBoards godoc
// @Summary List bingo boards
// @Description Lists the current player's bingo boards, cells, mission progress, and line reward state.
// @Tags Bingo
// @Produce json
// @Security AuthCookieAuth
// @Param category query string false "Board category" Enums(daily,persistent,event)
// @Success 200 {object} apimodel.BingoBoardListResponse
// @Failure 401 {object} httpx.ProblemDetails
// @Failure 500 {object} httpx.ProblemDetails
// @Router /bingo/boards [get]
func (h *Handler) ListBoards(w http.ResponseWriter, _ *http.Request) {
	httpx.WriteJSON(w, http.StatusOK, apimodel.BingoBoardListResponse{
		Boards: []apimodel.BingoBoardSummary{
			{
				BoardID:  "board_day_1",
				Title:    "Day 1 Bingo",
				Category: "daily",
				Status:   "active",
				Cells: []apimodel.BingoCell{
					{
						Row:    0,
						Column: 0,
						Mission: apimodel.MissionSummary{
							ID:          "mission_daily_match_3",
							Tab:         "daily",
							Title:       "Answer three match questions",
							Description: "Complete three Knowledge King questions today.",
							Status:      "claimable",
							Progress:    apimodel.MissionProgress{Current: 3, Target: 3},
							Rewards:     apimodel.Reward{OpenPower: 120},
						},
					},
				},
				Rewards: []apimodel.BingoReward{
					{
						LineRewardID: "line_reward_row_0",
						LineKey:      "row_0",
						Status:       "claimable",
						Reward:       apimodel.Reward{OpenPower: 200},
					},
				},
			},
		},
	})
}
