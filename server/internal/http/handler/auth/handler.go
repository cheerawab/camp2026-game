package auth

import "github.com/go-chi/chi/v5"

type Handler struct{}

func New() *Handler {
	return &Handler{}
}

func (h *Handler) RegisterRoutes(api chi.Router) {
	api.Post("/auth/login", h.Login)
	api.Post("/auth/logout", h.Logout)
}
