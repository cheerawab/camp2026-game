package me

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/sitcon-tw/camp2026-game/internal/content"
	"github.com/sitcon-tw/camp2026-game/internal/http/authctx"
	"github.com/sitcon-tw/camp2026-game/internal/http/httpx"
	mongomodel "github.com/sitcon-tw/camp2026-game/internal/mongodb/model"
)

type Dependencies struct {
	Content *content.Store
	MongoDB *mongo.Database
}

type Handler struct {
	content *content.Store
	db      *mongo.Database
}

func New(dep Dependencies) *Handler {
	return &Handler{
		content: dep.Content,
		db:      dep.MongoDB,
	}
}

func (h *Handler) RegisterRoutes(api chi.Router) {
	api.Get("/me/home", h.Home)
	api.Get("/me/status", h.Status)
	api.Get("/me/sitones", h.ListSitones)
	api.Get("/me/sitone-loadout", h.SitoneLoadout)
	api.Put("/me/sitone-loadout", h.UpdateSitoneLoadout)
	api.Get("/me/items", h.ListItems)
	api.Get("/me/qrcode", h.QRCode)
	api.Get("/me/matches", h.ListCompletedMatches)
}

func currentPlayer(w http.ResponseWriter, r *http.Request) (mongomodel.Player, bool) {
	player, ok := authctx.PlayerFromContext(r.Context())
	if !ok {
		httpx.WriteProblem(w, r, httpx.NewError(http.StatusUnauthorized, "authentication required"))
		return mongomodel.Player{}, false
	}
	return player, true
}
