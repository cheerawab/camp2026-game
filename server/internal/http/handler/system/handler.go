package system

import (
	"context"

	"github.com/go-chi/chi/v5"
)

type Dependencies struct {
	ReadinessCheck func(context.Context) error
}

type Handler struct {
	readinessCheck func(context.Context) error
}

func New(dep Dependencies) *Handler {
	return &Handler{
		readinessCheck: dep.ReadinessCheck,
	}
}

func (h *Handler) RegisterRoutes(api chi.Router) {
	api.Get("/healthz", h.Health)
}
