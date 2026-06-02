package catalog

import (
	"net/http"

	"github.com/sitcon-tw/camp2026-game/internal/content"
	"github.com/sitcon-tw/camp2026-game/internal/http/httpx"
)

type SitoneListResponse struct {
	Sitones []SitoneResponse `json:"sitones"`
}

type SitoneResponse struct {
	ID          string `json:"id" example:"sitone-engineering"`
	Name        string `json:"name" example:"工程型小石"`
	Type        string `json:"type" example:"engineering"`
	Rarity      string `json:"rarity" example:"base"`
	Style       string `json:"style" example:"default"`
	Description string `json:"description" example:"修 bug、分享解法、完成技術任務。"`
}

// ListSitones godoc
// @Summary List sitone catalog
// @Description Lists all collectible sitone definitions.
// @Tags catalog
// @Produce json
// @Success 200 {object} SitoneListResponse
// @Failure 503 {object} httpx.ProblemDetails
// @Router /catalog/sitones [get]
func (h *Handler) ListSitones(w http.ResponseWriter, r *http.Request) {
	if h.content == nil {
		httpx.WriteProblem(w, r, httpx.ServiceUnavailable("content store is unavailable"))
		return
	}

	httpx.WriteJSON(w, http.StatusOK, SitoneListResponse{
		Sitones: mapSitones(h.content.ListSitones()),
	})
}

func mapSitones(sitones []content.Sitone) []SitoneResponse {
	if len(sitones) == 0 {
		return nil
	}

	out := make([]SitoneResponse, 0, len(sitones))
	for _, sitone := range sitones {
		out = append(out, SitoneResponse{
			ID:          sitone.ID,
			Name:        sitone.Name,
			Type:        sitone.Type,
			Rarity:      sitone.Rarity,
			Style:       sitone.Style,
			Description: sitone.Description,
		})
	}
	return out
}
