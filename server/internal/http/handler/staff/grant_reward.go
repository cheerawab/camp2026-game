package staff

import (
	"net/http"

	"github.com/sitcon-tw/camp2026-game/internal/http/apimodel"
	"github.com/sitcon-tw/camp2026-game/internal/http/httpx"
)

// GrantReward godoc
// @Summary Grant staff reward
// @Description Lets staff grant open power, sitones, or items to a player resolved from a QRCode token.
// @Tags Staff
// @Accept json
// @Produce json
// @Security AuthCookieAuth
// @Param request body apimodel.StaffGrantRewardRequest true "Staff reward grant request"
// @Success 201 {object} apimodel.StaffGrantRewardResponse
// @Failure 400 {object} httpx.ProblemDetails
// @Failure 401 {object} httpx.ProblemDetails
// @Failure 422 {object} httpx.ProblemDetails
// @Failure 501 {object} httpx.ProblemDetails
// @Router /staff/rewards [post]
func (h *Handler) GrantReward(w http.ResponseWriter, r *http.Request) {
	var body apimodel.StaffGrantRewardRequest
	if err := httpx.DecodeJSON(r, &body); err != nil {
		httpx.WriteProblem(w, r, err)
		return
	}
	if err := httpx.ValidateStruct(body); err != nil {
		httpx.WriteProblem(w, r, err)
		return
	}

	httpx.WriteProblem(w, r, httpx.NotImplemented("staff reward grants are not implemented yet"))
}
