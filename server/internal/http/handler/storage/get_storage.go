package storage

import (
	"net/http"

	"github.com/sitcon-tw/camp2026-game/internal/http/apimodel"
	"github.com/sitcon-tw/camp2026-game/internal/http/httpx"
)

// GetStorage godoc
// @Summary Get storage summary
// @Description Returns the current player's collection summary, including owned sitones, items, and craftable recipe count.
// @Tags Storage
// @Produce json
// @Security AuthCookieAuth
// @Success 200 {object} apimodel.StorageSummaryResponse
// @Failure 401 {object} httpx.ProblemDetails
// @Failure 500 {object} httpx.ProblemDetails
// @Router /storage [get]
func (h *Handler) GetStorage(w http.ResponseWriter, _ *http.Request) {
	httpx.WriteJSON(w, http.StatusOK, apimodel.StorageSummaryResponse{
		Sitones: []apimodel.SitoneSummary{
			{
				ID:           "sitone_01HR9Z7E2Z2VJ2QZ4P4Z",
				DefinitionID: "sitone-engineering",
				Name:         "Engineering Sitone",
				Type:         "engineering",
				Rarity:       "rare",
				Style:        "default",
			},
		},
		Items: []apimodel.ItemSummary{
			{
				ID:           "pit_01HR9Z7E2Z2VJ2QZ4P4Z",
				DefinitionID: "item-camp-sticker",
				Name:         "Camp Sticker",
				Quantity:     3,
			},
		},
		CraftableRecipes: 1,
	})
}
