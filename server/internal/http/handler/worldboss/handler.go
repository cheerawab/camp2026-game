package worldboss

import "github.com/go-chi/chi/v5"

type Handler struct{}

func New() *Handler {
	return &Handler{}
}

func (h *Handler) RegisterRoutes(api chi.Router) {
	api.Get("/world-bosses", h.ListBosses)
	api.Get("/world-bosses/{bossID}", h.GetBoss)
	api.Post("/world-bosses/{bossID}/matches", h.CreateMatch)
}
