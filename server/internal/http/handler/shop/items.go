package shop

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/sitcon-tw/camp2026-game/internal/content"
	"github.com/sitcon-tw/camp2026-game/internal/http/httpx"
)

// ListItems godoc
// @Summary List shop items
// @Description Lists purchasable and enabled item definitions.
// @Tags shop
// @Produce json
// @Security AuthCookieAuth
// @Success 200 {object} ItemListResponse
// @Failure 401 {object} httpx.ProblemDetails
// @Failure 503 {object} httpx.ProblemDetails
// @Router /shop/items [get]
func (h *Handler) ListItems(w http.ResponseWriter, r *http.Request) {
	if _, ok := currentPlayer(w, r); !ok {
		return
	}
	if !h.requireContent(w, r) {
		return
	}

	httpx.WriteJSON(w, http.StatusOK, ItemListResponse{
		Items: shopItemResponses(shopItems(h.content)),
	})
}

// GetItem godoc
// @Summary Get shop item
// @Description Returns one purchasable and enabled item definition.
// @Tags shop
// @Produce json
// @Security AuthCookieAuth
// @Param itemID path string true "Item ID"
// @Success 200 {object} ItemDetailResponse
// @Failure 401 {object} httpx.ProblemDetails
// @Failure 404 {object} httpx.ProblemDetails
// @Failure 503 {object} httpx.ProblemDetails
// @Router /shop/items/{itemID} [get]
func (h *Handler) GetItem(w http.ResponseWriter, r *http.Request) {
	if _, ok := currentPlayer(w, r); !ok {
		return
	}
	if !h.requireContent(w, r) {
		return
	}

	itemID := chi.URLParam(r, "itemID")
	item, ok := shopItemByID(h.content, itemID)
	if !ok {
		httpx.WriteProblem(w, r, httpx.NotFound("shop item not found"))
		return
	}

	httpx.WriteJSON(w, http.StatusOK, ItemDetailResponse{
		Item: shopItemResponse(item),
	})
}

func shopItems(store *content.Store) []content.Item {
	items := store.ListItems()
	out := make([]content.Item, 0, len(items))
	for _, item := range items {
		if item.Purchasable && item.Enabled {
			out = append(out, item)
		}
	}
	return out
}

func shopItemByID(store *content.Store, itemID string) (content.Item, bool) {
	item, ok := store.GetItem(itemID)
	if !ok || !item.Purchasable || !item.Enabled {
		return content.Item{}, false
	}
	return item, true
}

func shopItemResponses(items []content.Item) []ShopItemResponse {
	if len(items) == 0 {
		return nil
	}

	out := make([]ShopItemResponse, 0, len(items))
	for _, item := range items {
		out = append(out, shopItemResponse(item))
	}
	return out
}

func shopItemResponse(item content.Item) ShopItemResponse {
	return ShopItemResponse{
		ID:             item.ID,
		Name:           item.Name,
		Type:           item.Type,
		Rarity:         item.Rarity,
		Description:    item.Description,
		PriceOpenPower: item.PriceOpenPower,
	}
}
