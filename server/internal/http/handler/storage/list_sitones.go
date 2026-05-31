package storage

import (
	"net/http"

	"github.com/sitcon-tw/camp2026-game/internal/http/apimodel"
	"github.com/sitcon-tw/camp2026-game/internal/http/httpx"
)

// ListSitones godoc
// @Summary List owned sitones
// @Description Lists sitones owned by the current player. Sitones can be assigned to loadout or used in crafting.
// @Tags Storage
// @Produce json
// @Security AuthCookieAuth
// @Success 200 {object} apimodel.SitoneListResponse
// @Failure 401 {object} httpx.ProblemDetails
// @Failure 500 {object} httpx.ProblemDetails
// @Router /storage/sitones [get]
func (h *Handler) ListSitones(w http.ResponseWriter, _ *http.Request) {
	httpx.WriteJSON(w, http.StatusOK, apimodel.SitoneListResponse{
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
	})
}
