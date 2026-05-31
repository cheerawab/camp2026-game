package bingo

import (
	"net/http"

	"github.com/sitcon-tw/camp2026-game/internal/http/httpx"
)

// ClaimLineReward godoc
// @Summary Claim bingo line reward
// @Description Claims a bingo line reward after the required line is completed.
// @Tags Bingo
// @Produce json
// @Security AuthCookieAuth
// @Param lineRewardID path string true "Line reward ID"
// @Success 200 {object} apimodel.BingoLineRewardClaimResponse
// @Failure 400 {object} httpx.ProblemDetails
// @Failure 401 {object} httpx.ProblemDetails
// @Failure 404 {object} httpx.ProblemDetails
// @Failure 501 {object} httpx.ProblemDetails
// @Router /bingo/line-rewards/{lineRewardID}/claim [post]
func (h *Handler) ClaimLineReward(w http.ResponseWriter, r *http.Request) {
	httpx.WriteProblem(w, r, httpx.NotImplemented("bingo line reward claim is not implemented yet"))
}
