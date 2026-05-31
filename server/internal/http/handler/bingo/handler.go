package bingo

import "github.com/go-chi/chi/v5"

type Handler struct{}

func New() *Handler {
	return &Handler{}
}

func (h *Handler) RegisterRoutes(api chi.Router) {
	api.Get("/bingo/boards", h.ListBoards)
	api.Post("/bingo/missions/{missionID}/complete", h.CompleteMission)
	api.Post("/bingo/line-rewards/{lineRewardID}/claim", h.ClaimLineReward)
}
