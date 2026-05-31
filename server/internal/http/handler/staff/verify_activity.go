package staff

import (
	"net/http"

	"github.com/sitcon-tw/camp2026-game/internal/http/apimodel"
	"github.com/sitcon-tw/camp2026-game/internal/http/httpx"
)

// VerifyActivity godoc
// @Summary Verify staff activity
// @Description Lets staff verify a player activity from a QRCode token and optionally advance a mission.
// @Tags Staff
// @Accept json
// @Produce json
// @Security AuthCookieAuth
// @Param request body apimodel.StaffActivityVerificationRequest true "Staff activity verification request"
// @Success 201 {object} apimodel.StaffActivityVerificationResponse
// @Failure 400 {object} httpx.ProblemDetails
// @Failure 401 {object} httpx.ProblemDetails
// @Failure 422 {object} httpx.ProblemDetails
// @Failure 501 {object} httpx.ProblemDetails
// @Router /staff/activity-verifications [post]
func (h *Handler) VerifyActivity(w http.ResponseWriter, r *http.Request) {
	var body apimodel.StaffActivityVerificationRequest
	if err := httpx.DecodeJSON(r, &body); err != nil {
		httpx.WriteProblem(w, r, err)
		return
	}
	if err := httpx.ValidateStruct(body); err != nil {
		httpx.WriteProblem(w, r, err)
		return
	}

	httpx.WriteProblem(w, r, httpx.NotImplemented("staff activity verification is not implemented yet"))
}
