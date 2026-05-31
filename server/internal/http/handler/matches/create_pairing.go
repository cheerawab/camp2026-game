package matches

import (
	"net/http"

	"github.com/sitcon-tw/camp2026-game/internal/http/apimodel"
	"github.com/sitcon-tw/camp2026-game/internal/http/httpx"
)

// CreatePairing godoc
// @Summary Create match QRCode pairing
// @Description Creates a match pairing after scanning another player's QRCode.
// @Tags Matches
// @Accept json
// @Produce json
// @Security AuthCookieAuth
// @Param request body apimodel.MatchPairingRequest true "QRCode pairing request"
// @Success 201 {object} apimodel.MatchPairingResponse
// @Failure 400 {object} httpx.ProblemDetails
// @Failure 401 {object} httpx.ProblemDetails
// @Failure 422 {object} httpx.ProblemDetails
// @Failure 501 {object} httpx.ProblemDetails
// @Router /match-pairings [post]
func (h *Handler) CreatePairing(w http.ResponseWriter, r *http.Request) {
	var body apimodel.MatchPairingRequest
	if err := httpx.DecodeJSON(r, &body); err != nil {
		httpx.WriteProblem(w, r, err)
		return
	}
	if err := httpx.ValidateStruct(body); err != nil {
		httpx.WriteProblem(w, r, err)
		return
	}

	httpx.WriteProblem(w, r, httpx.NotImplemented("match pairing is not implemented yet"))
}
