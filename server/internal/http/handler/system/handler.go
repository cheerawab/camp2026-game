package system

import (
	"context"

	"github.com/go-chi/chi/v5"

	"github.com/sitcon-tw/camp2026-game/internal/postgres/sqlc"
)

type Dependencies struct {
	ReadinessCheck func(context.Context) error
	Queries        *sqlc.Queries
}

type Handler struct {
	readinessCheck func(context.Context) error
	queries        *sqlc.Queries
}

func New(dep Dependencies) *Handler {
	return &Handler{
		readinessCheck: dep.ReadinessCheck,
		queries:        dep.Queries,
	}
}

func (h *Handler) RegisterRoutes(api chi.Router) {
	api.Get("/", h.Root)
	api.Get("/healthz", h.Health)
	api.Get("/readyz", h.Ready)
	api.Get("/ping", h.Ping)
}
