package catalog

import (
	"github.com/go-chi/chi/v5"

	"github.com/sitcon-tw/camp2026-game/internal/content"
)

type Dependencies struct {
	Content *content.Store
}

type Handler struct {
	content *content.Store
}

func New(dep Dependencies) *Handler {
	return &Handler{
		content: dep.Content,
	}
}

func (h *Handler) RegisterRoutes(api chi.Router) {
	api.Get("/catalog/sitones", h.ListSitones)
}
