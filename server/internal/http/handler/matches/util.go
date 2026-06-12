package matches

import (
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
	"strings"
	"time"

	"github.com/sitcon-tw/camp2026-game/internal/content"
	mongomodel "github.com/sitcon-tw/camp2026-game/internal/mongodb/model"
)

const (
	idAlphabet   = "abcdefghijklmnopqrstuvwxyz0123456789"
	codeAlphabet = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
)

func newID(prefix string) (string, error) {
	value, err := randomString(idAlphabet, 12)
	if err != nil {
		return "", err
	}
	return prefix + "_" + value, nil
}

func newMatchCode() (string, error) {
	return randomString(codeAlphabet, 6)
}

func randomString(alphabet string, length int) (string, error) {
	var out strings.Builder
	out.Grow(length)

	max := big.NewInt(int64(len(alphabet)))
	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", fmt.Errorf("generate random string: %w", err)
		}
		out.WriteByte(alphabet[n.Int64()])
	}
	return out.String(), nil
}

func playerIndex(match mongomodel.Match, playerID string) int {
	for i, player := range match.Players {
		if player.PlayerID == playerID {
			return i
		}
	}
	return -1
}

func isParticipant(match mongomodel.Match, playerID string) bool {
	return playerIndex(match, playerID) >= 0
}

func allPlayersReady(match mongomodel.Match) bool {
	if len(match.Players) != 2 {
		return false
	}
	for _, player := range match.Players {
		if !player.Ready || len(player.SitoneIDs) == 0 {
			return false
		}
	}
	return true
}

func activeMatchPhase(match mongomodel.Match) string {
	if match.Status != mongomodel.MatchStatusActive {
		return ""
	}
	if match.Phase == mongomodel.MatchPhaseRevealing {
		return mongomodel.MatchPhaseRevealing
	}
	return mongomodel.MatchPhaseAnswering
}

func questionResponse(question content.QuizQuestion) MatchQuestionResponse {
	return MatchQuestionResponse{
		QuestionID: question.ID,
		Prompt:     question.Prompt,
		ChoiceA:    question.ChoiceA,
		ChoiceB:    question.ChoiceB,
		ChoiceC:    question.ChoiceC,
		ChoiceD:    question.ChoiceD,
	}
}

func scoreAnswer(question content.QuizQuestion, choice string, answeredAt time.Time, roundEndsAt time.Time) (bool, int) {
	correct := strings.EqualFold(choice, question.CorrectChoice)
	if !correct {
		return false, 0
	}

	remainingSeconds := int(math.Floor(roundEndsAt.Sub(answeredAt).Seconds()))
	if remainingSeconds < 0 {
		remainingSeconds = 0
	}
	return true, 100 + remainingSeconds*5
}

func maxScorePerQuestion() int {
	return 100 + roundDuration*5
}

func maxScoreThroughCurrentQuestion(match mongomodel.Match) int {
	if len(match.QuestionIDs) == 0 || match.Status == mongomodel.MatchStatusWaiting {
		return 0
	}

	questionCount := match.CurrentQuestionIndex + 1
	if match.Status == mongomodel.MatchStatusCompleted {
		questionCount = len(match.QuestionIDs)
	}
	if questionCount < 0 {
		questionCount = 0
	}
	if questionCount > len(match.QuestionIDs) {
		questionCount = len(match.QuestionIDs)
	}

	return questionCount * maxScorePerQuestion()
}

func openPowerReward(score int) int {
	return score/10 + 20
}

func matchHasClearWinner(match mongomodel.Match) bool {
	if len(match.Players) < 2 {
		return false
	}

	topScore := match.Players[0].Score
	topCount := 1
	for _, player := range match.Players[1:] {
		switch {
		case player.Score > topScore:
			topScore = player.Score
			topCount = 1
		case player.Score == topScore:
			topCount++
		}
	}
	return topCount == 1
}

func matchOpenPowerReward(match mongomodel.Match, player mongomodel.MatchPlayer) int {
	if !matchHasClearWinner(match) {
		return 0
	}
	return openPowerReward(player.Score)
}

func matchRewardRecordID(matchID, playerID string) string {
	return "open_power_reward_" + matchID + "_" + playerID
}

func matchRewardSource(matchID, playerID string) string {
	return "quiz_match:" + matchID + ":player:" + playerID
}

func timePtr(value time.Time) *time.Time {
	if value.IsZero() {
		return nil
	}
	return &value
}
