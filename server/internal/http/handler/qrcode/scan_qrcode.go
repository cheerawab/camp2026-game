package qrcode

import (
	"net/http"

	"github.com/sitcon-tw/camp2026-game/internal/http/apimodel"
	"github.com/sitcon-tw/camp2026-game/internal/http/httpx"
)

// ScanQRCode godoc
// @Summary Scan QRCode
// @Description Resolves a player QRCode and returns available user actions for match pairing or staff verification.
// @Tags QRCode
// @Accept json
// @Produce json
// @Security AuthCookieAuth
// @Param request body apimodel.QRCodeScanRequest true "QRCode scan request"
// @Success 200 {object} apimodel.QRCodeScanResponse
// @Failure 400 {object} httpx.ProblemDetails
// @Failure 401 {object} httpx.ProblemDetails
// @Failure 422 {object} httpx.ProblemDetails
// @Failure 501 {object} httpx.ProblemDetails
// @Router /qrcode/scans [post]
func (h *Handler) ScanQRCode(w http.ResponseWriter, r *http.Request) {
	var body apimodel.QRCodeScanRequest
	if err := httpx.DecodeJSON(r, &body); err != nil {
		httpx.WriteProblem(w, r, err)
		return
	}
	if err := httpx.ValidateStruct(body); err != nil {
		httpx.WriteProblem(w, r, err)
		return
	}

	httpx.WriteProblem(w, r, httpx.NotImplemented("qrcode scanning is not implemented yet"))
}
