package storage

import (
	"net/http"

	"github.com/sitcon-tw/camp2026-game/internal/http/apimodel"
	"github.com/sitcon-tw/camp2026-game/internal/http/httpx"
)

// ListRecipes godoc
// @Summary List visible crafting recipes
// @Description Lists crafting recipes visible to the current player, including whether each recipe is currently craftable.
// @Tags Storage
// @Produce json
// @Security AuthCookieAuth
// @Success 200 {object} apimodel.RecipeListResponse
// @Failure 401 {object} httpx.ProblemDetails
// @Failure 500 {object} httpx.ProblemDetails
// @Router /storage/recipes [get]
func (h *Handler) ListRecipes(w http.ResponseWriter, _ *http.Request) {
	httpx.WriteJSON(w, http.StatusOK, apimodel.RecipeListResponse{
		Recipes: []apimodel.RecipeSummary{
			{
				RecipeID:             "recipe_engineering_skin",
				Name:                 "Engineering Skin",
				RequiredSitoneType:   "engineering",
				RequiredItemID:       "item-camp-sticker",
				RequiredItemQuantity: 1,
				OutputDefinitionID:   "sitone-engineering-skin",
				Unlocked:             true,
				Craftable:            true,
				AcquisitionHint:      "Complete engineering missions.",
			},
		},
	})
}
