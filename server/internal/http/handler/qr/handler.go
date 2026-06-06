package qr

import (
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/go-chi/chi/v5"
)

type Dependencies struct {
	MongoDB *mongo.Database
}

type Handler struct {
	db *mongo.Database
}

func New(dep Dependencies) *Handler {
	return &Handler{db: dep.MongoDB}
}

func (h *Handler) RegisterRoutes(api chi.Router) {
	api.Post("/qr/resolve", h.Resolve)
}
