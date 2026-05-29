package example

import (
	"net/http"

	"github.com/sitcon-tw/camp2026-game/internal/http/httpx"
)

type ValidationExampleRequest struct {
	Players []PlayerInput `json:"players" validate:"required,min=1,max=20,dive"`
}

type PlayerInput struct {
	DisplayName        string `json:"displayName" validate:"required,max=40"`
	TeamNumber         int    `json:"teamNumber" validate:"required,min=1,max=20"`
	FavoritePebbleType string `json:"favoritePebbleType" validate:"required,oneof=exploration inspiration resonance engineering entertainment"`
}

type ValidationExampleResponse struct {
	Message string        `json:"message" example:"validation example accepted"`
	Players []PlayerInput `json:"players"`
}

// Validation godoc
// @Summary Validation example
// @Description Demonstrates JSON decode errors, validator semantic errors, and array-level error locations.
// @Tags examples
// @Accept json
// @Produce json
// @Param request body ValidationExampleRequest true "Validation example request"
// @Success 201 {object} ValidationExampleResponse
// @Failure 400 {object} httpx.ProblemDetails
// @Failure 422 {object} httpx.ProblemDetails
// @Router /examples/validation [post]
func (h *Handler) Validation(w http.ResponseWriter, r *http.Request) {
	var body ValidationExampleRequest
	if err := httpx.DecodeJSON(r, &body); err != nil {
		httpx.WriteProblem(w, r, err)
		return
	}

	normalizeValidationExampleRequest(&body)
	if err := httpx.ValidateStruct(body); err != nil {
		httpx.WriteProblem(w, r, err)
		return
	}

	httpx.WriteJSON(w, http.StatusCreated, ValidationExampleResponse{
		Message: "validation example accepted",
		Players: body.Players,
	})
}
