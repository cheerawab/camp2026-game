package auth

import (
	"net/http"

	"github.com/sitcon-tw/camp2026-game/internal/http/apimodel"
	"github.com/sitcon-tw/camp2026-game/internal/http/httpx"
)

// Login godoc
// @Summary User login with auth token
// @Description User-facing endpoint. Validates the stable auth token from the issued game URL and writes it to the camp2026_auth cookie.
// @Tags User Authentication
// @Accept json
// @Produce json
// @Param request body apimodel.AuthLoginRequest true "Auth login request"
// @Success 200 {object} apimodel.AuthLoginResponse
// @Header 200 {string} Set-Cookie "camp2026_auth=<auth-token>; Path=/; HttpOnly; Secure; SameSite=Lax"
// @Failure 400 {object} httpx.ProblemDetails
// @Failure 401 {object} httpx.ProblemDetails
// @Failure 422 {object} httpx.ProblemDetails
// @Failure 501 {object} httpx.ProblemDetails
// @Router /auth/login [post]
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var body apimodel.AuthLoginRequest
	if err := httpx.DecodeJSON(r, &body); err != nil {
		httpx.WriteProblem(w, r, err)
		return
	}
	if err := httpx.ValidateStruct(body); err != nil {
		httpx.WriteProblem(w, r, err)
		return
	}

	httpx.WriteProblem(w, r, httpx.NotImplemented("auth login is not implemented yet"))
}
