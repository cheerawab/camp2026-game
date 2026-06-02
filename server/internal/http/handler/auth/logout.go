package auth

import (
	"net/http"

	"github.com/sitcon-tw/camp2026-game/internal/http/authctx"
)

// Logout godoc
// @Summary User logout
// @Description Clears the camp2026_auth cookie for the current browser session.
// @Tags User Authentication
// @Produce json
// @Security AuthCookieAuth
// @Success 204
// @Header 204 {string} Set-Cookie "camp2026_auth=; Path=/; HttpOnly; Secure; SameSite=Lax; Max-Age=0"
// @Failure 401 {object} httpx.ProblemDetails
// @Router /auth/logout [post]
func (h *Handler) Logout(w http.ResponseWriter, _ *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     authctx.CookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
	w.WriteHeader(http.StatusNoContent)
}
