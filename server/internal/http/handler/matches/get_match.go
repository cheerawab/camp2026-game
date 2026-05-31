package matches

import (
	"net/http"

	"github.com/sitcon-tw/camp2026-game/internal/http/apimodel"
	"github.com/sitcon-tw/camp2026-game/internal/http/httpx"
)

// GetMatch godoc
// @Summary Get match state
// @Description Returns match status, participants, current question, or completed result for the current player.
// @Tags Matches
// @Produce json
// @Security AuthCookieAuth
// @Param matchID path string true "Match ID"
// @Success 200 {object} apimodel.MatchResponse
// @Failure 401 {object} httpx.ProblemDetails
// @Failure 404 {object} httpx.ProblemDetails
// @Failure 500 {object} httpx.ProblemDetails
// @Router /matches/{matchID} [get]
func (h *Handler) GetMatch(w http.ResponseWriter, _ *http.Request) {
	httpx.WriteJSON(w, http.StatusOK, apimodel.MatchResponse{
		MatchID: "match_01HR9Z7E2Z2VJ2QZ4P4Z",
		Mode:    "qr_duel",
		Status:  "answering",
		Player: apimodel.MatchParticipant{
			PlayerID:    "player_01HR9Z7E2Z2VJ2QZ4P4Z",
			DisplayName: "Alice",
			HP:          100,
			OpenPower:   1280,
		},
		Opponent: apimodel.MatchParticipant{
			PlayerID:    "player_01HR9Z7E2Z2VJ2QZ4P4Y",
			DisplayName: "Bob",
			HP:          90,
			OpenPower:   1100,
		},
		Question: apimodel.MatchQuestion{
			QuestionID: "question_001",
			Prompt:     "Which license is commonly used for open source projects?",
			Choices: []apimodel.MatchChoice{
				{ChoiceID: "A", Text: "MIT License"},
				{ChoiceID: "B", Text: "Private NDA"},
			},
			TimeLimit: 30,
		},
	})
}
