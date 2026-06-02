package catalog

import (
	"net/http"

	"github.com/sitcon-tw/camp2026-game/internal/content"
	"github.com/sitcon-tw/camp2026-game/internal/http/httpx"
)

type ItemListResponse struct {
	Items []ItemResponse `json:"items"`
}

type ItemResponse struct {
	ID          string `json:"id" example:"item-crafting-fragment"`
	Name        string `json:"name" example:"合成碎片"`
	Type        string `json:"type" example:"material"`
	Rarity      string `json:"rarity" example:"common"`
	Description string `json:"description" example:"小石造型合成使用的基礎素材。"`
}

// ListItems godoc
// @Summary List item catalog
// @Description Lists all collectible item definitions.
// @Tags catalog
// @Produce json
// @Success 200 {object} ItemListResponse
// @Failure 503 {object} httpx.ProblemDetails
// @Router /catalog/items [get]
func (h *Handler) ListItems(w http.ResponseWriter, r *http.Request) {
	if h.content == nil {
		httpx.WriteProblem(w, r, httpx.ServiceUnavailable("content store is unavailable"))
		return
	}

	httpx.WriteJSON(w, http.StatusOK, ItemListResponse{
		Items: mapItems(h.content.ListItems()),
	})
}

func mapItems(items []content.Item) []ItemResponse {
	if len(items) == 0 {
		return nil
	}

	out := make([]ItemResponse, 0, len(items))
	for _, item := range items {
		out = append(out, ItemResponse{
			ID:          item.ID,
			Name:        item.Name,
			Type:        item.Type,
			Rarity:      item.Rarity,
			Description: item.Description,
		})
	}
	return out
}
