package storage

import "github.com/go-chi/chi/v5"

type Handler struct{}

func New() *Handler {
	return &Handler{}
}

func (h *Handler) RegisterRoutes(api chi.Router) {
	api.Get("/storage", h.GetStorage)
	api.Get("/storage/sitones", h.ListSitones)
	api.Get("/storage/items", h.ListItems)
	api.Get("/storage/recipes", h.ListRecipes)
	api.Post("/storage/crafting", h.CraftSitone)
}
