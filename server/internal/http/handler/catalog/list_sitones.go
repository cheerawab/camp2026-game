package catalog

import (
	"net/http"

	"github.com/sitcon-tw/camp2026-game/internal/http/apimodel"
	"github.com/sitcon-tw/camp2026-game/internal/http/httpx"
)

// ListSitones godoc
// @Summary List sitone catalog
// @Description Lists all collectible sitone definitions and acquisition hints.
// @Tags Catalog
// @Produce json
// @Security AuthCookieAuth
// @Success 200 {object} apimodel.CatalogSitoneListResponse
// @Failure 401 {object} httpx.ProblemDetails
// @Failure 500 {object} httpx.ProblemDetails
// @Router /catalog/sitones [get]
func (h *Handler) ListSitones(w http.ResponseWriter, _ *http.Request) {
	httpx.WriteJSON(w, http.StatusOK, apimodel.CatalogSitoneListResponse{
		Sitones: []apimodel.CatalogSitoneSummary{
			{
				DefinitionID:    "sitone-engineering",
				Name:            "Engineering Sitone",
				Type:            "engineering",
				Rarity:          "rare",
				AcquisitionHint: "Complete engineering missions.",
			},
		},
	})
}
