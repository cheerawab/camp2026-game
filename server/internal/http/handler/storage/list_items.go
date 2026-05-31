package storage

import (
	"net/http"

	"github.com/sitcon-tw/camp2026-game/internal/http/apimodel"
	"github.com/sitcon-tw/camp2026-game/internal/http/httpx"
)

// ListItems godoc
// @Summary List owned items
// @Description Lists items owned by the current player. Items are mainly used for crafting and event collection.
// @Tags Storage
// @Produce json
// @Security AuthCookieAuth
// @Success 200 {object} apimodel.ItemListResponse
// @Failure 401 {object} httpx.ProblemDetails
// @Failure 500 {object} httpx.ProblemDetails
// @Router /storage/items [get]
func (h *Handler) ListItems(w http.ResponseWriter, _ *http.Request) {
	httpx.WriteJSON(w, http.StatusOK, apimodel.ItemListResponse{
		Items: []apimodel.ItemSummary{
			{
				ID:           "pit_01HR9Z7E2Z2VJ2QZ4P4Z",
				DefinitionID: "item-camp-sticker",
				Name:         "Camp Sticker",
				Quantity:     3,
			},
		},
	})
}
