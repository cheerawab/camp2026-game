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
		if !player.Ready {
			return false
		}
	}
	return true
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

func openPowerReward(score int) int {
	return score/10 + 20
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
