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
	ID          string `json:"id" example:"item_adventure_backpack"`
	Name        string `json:"name" example:"冒險背包"`
	Type        string `json:"type" example:"material"`
	Rarity      string `json:"rarity" example:"common"`
	Description string `json:"description" example:"冒險背包，可用於小石合成。"`
	IconPath    string `json:"iconPath,omitempty" example:"/game-icons/items/item_adventure_backpack.png"`
	Source      string `json:"source,omitempty" example:"shop"`
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
			IconPath:    item.IconPath,
			Source:      item.Source,
		})
	}
	return out
}
