package matches

import (
	"net/http"

	"github.com/sitcon-tw/camp2026-game/internal/http/apimodel"
	"github.com/sitcon-tw/camp2026-game/internal/http/httpx"
)

// ListMatches godoc
// @Summary List match history
// @Description Lists the current player's match history. Results can be filtered by mode and status.
// @Tags Matches
// @Produce json
// @Security AuthCookieAuth
// @Param mode query string false "Match mode" Enums(qr_duel,offline_duel,world_boss)
// @Param status query string false "Match status" Enums(pairing,answering,completed,cancelled)
// @Success 200 {object} apimodel.MatchListResponse
// @Failure 401 {object} httpx.ProblemDetails
// @Failure 500 {object} httpx.ProblemDetails
// @Router /matches [get]
func (h *Handler) ListMatches(w http.ResponseWriter, _ *http.Request) {
	httpx.WriteJSON(w, http.StatusOK, apimodel.MatchListResponse{
		Matches: []apimodel.MatchSummary{
			{
				MatchID:         "match_01HR9Z7E2Z2VJ2QZ4P4Z",
				Mode:            "qr_duel",
				Status:          "completed",
				OpponentName:    "Bob",
				PlayerScore:     320,
				OpponentScore:   280,
				CompletedAt:     "2026-07-24T10:35:00+08:00",
				OpenPowerGained: 80,
			},
		},
	})
}
