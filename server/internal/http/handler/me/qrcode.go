package me

import (
	"net/http"

	"github.com/sitcon-tw/camp2026-game/internal/http/httpx"
)

// QRCode godoc
// @Summary Get current player QR code token
// @Description Returns the authenticated player's QR code token for client-side QR rendering.
// @Tags me
// @Produce json
// @Security AuthCookieAuth
// @Success 200 {object} QRCodeResponse
// @Failure 401 {object} httpx.ProblemDetails
// @Failure 500 {object} httpx.ProblemDetails
// @Router /me/qrcode [get]
func (h *Handler) QRCode(w http.ResponseWriter, r *http.Request) {
	player, ok := currentPlayer(w, r)
	if !ok {
		return
	}
	if player.QRCodeToken == "" {
		httpx.WriteProblem(w, r, httpx.NewError(http.StatusInternalServerError, "qrcode is unavailable"))
		return
	}

	httpx.WriteJSON(w, http.StatusOK, QRCodeResponse{
		QRCodeToken: player.QRCodeToken,
	})
}
