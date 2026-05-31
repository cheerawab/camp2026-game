package staff

import "github.com/go-chi/chi/v5"

type Handler struct{}

func New() *Handler {
	return &Handler{}
}

func (h *Handler) RegisterRoutes(api chi.Router) {
	api.Post("/staff/rewards", h.GrantReward)
	api.Post("/staff/activity-verifications", h.VerifyActivity)
}
