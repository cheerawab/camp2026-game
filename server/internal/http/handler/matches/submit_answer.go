package matches

import (
	"net/http"

	"github.com/sitcon-tw/camp2026-game/internal/http/apimodel"
	"github.com/sitcon-tw/camp2026-game/internal/http/httpx"
)

// SubmitAnswer godoc
// @Summary Submit match answer
// @Description Submits an answer, returns correctness and explanation, and advances match state.
// @Tags Matches
// @Accept json
// @Produce json
// @Security AuthCookieAuth
// @Param matchID path string true "Match ID"
// @Param request body apimodel.MatchAnswerSubmitRequest true "Match answer request"
// @Success 200 {object} apimodel.MatchAnswerSubmitResponse
// @Failure 400 {object} httpx.ProblemDetails
// @Failure 401 {object} httpx.ProblemDetails
// @Failure 422 {object} httpx.ProblemDetails
// @Failure 501 {object} httpx.ProblemDetails
// @Router /matches/{matchID}/answers [post]
func (h *Handler) SubmitAnswer(w http.ResponseWriter, r *http.Request) {
	var body apimodel.MatchAnswerSubmitRequest
	if err := httpx.DecodeJSON(r, &body); err != nil {
		httpx.WriteProblem(w, r, err)
		return
	}
	if err := httpx.ValidateStruct(body); err != nil {
		httpx.WriteProblem(w, r, err)
		return
	}

	httpx.WriteProblem(w, r, httpx.NotImplemented("match answer submission is not implemented yet"))
}
