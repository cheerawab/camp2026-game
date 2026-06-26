package matches

import (
	"context"
	"sort"
	"strings"

	"github.com/sitcon-tw/camp2026-game/internal/content"
	mongomodel "github.com/sitcon-tw/camp2026-game/internal/mongodb/model"
)

func (h *Handler) ensureCurrentRoundEliminations(ctx context.Context, match *mongomodel.Match) error {
	if match == nil || match.Status != mongomodel.MatchStatusActive {
		return nil
	}
	if match.CurrentQuestionIndex < 0 || match.CurrentQuestionIndex >= len(match.QuestionIDs) {
		return nil
	}

	questionID := match.QuestionIDs[match.CurrentQuestionIndex]
	question, ok := h.content.GetQuizQuestion(questionID)
	if !ok {
		return nil
	}

	existing := make(map[string]struct{}, len(match.EliminatedChoices))
	for _, eliminated := range match.EliminatedChoices {
		existing[eliminationKey(eliminated.QuestionID, eliminated.PlayerID)] = struct{}{}
	}

	for _, player := range match.Players {
		key := eliminationKey(questionID, player.PlayerID)
		if _, ok := existing[key]; ok {
			continue
		}

		effects, err := h.matchPlayerBattleEffects(ctx, player)
		if err != nil {
			return err
		}
		match.EliminatedChoices = append(match.EliminatedChoices, eliminatedChoicesForPlayer(match.ID, question, player, effects))
	}
	return nil
}

func eliminatedChoicesForPlayer(
	matchID string,
	question content.QuizQuestion,
	player mongomodel.MatchPlayer,
	effects battleEffects,
) mongomodel.MatchEliminatedChoice {
	record := mongomodel.MatchEliminatedChoice{
		QuestionID: question.ID,
		PlayerID:   player.PlayerID,
	}
	if effects.EliminateChancePercent <= 0 || effects.EliminateCount <= 0 {
		return record
	}
	if deterministicPercent("eliminate", matchID, question.ID, player.PlayerID) >= effects.EliminateChancePercent {
		return record
	}

	wrongChoices := wrongChoiceLabels(question)
	sort.Slice(wrongChoices, func(i, j int) bool {
		left := deterministicUint64("wrong-choice", matchID, question.ID, player.PlayerID, wrongChoices[i])
		right := deterministicUint64("wrong-choice", matchID, question.ID, player.PlayerID, wrongChoices[j])
		return left < right
	})

	count := effects.EliminateCount
	if count > len(wrongChoices) {
		count = len(wrongChoices)
	}
	record.Choices = cloneStrings(wrongChoices[:count])
	record.SourceSitoneNames = cloneStrings(effects.EliminateSourceNames)
	return record
}

func wrongChoiceLabels(question content.QuizQuestion) []string {
	correct := strings.ToUpper(strings.TrimSpace(question.CorrectChoice))
	choices := []string{"A", "B", "C", "D"}
	out := make([]string, 0, 3)
	for _, choice := range choices {
		if choice != correct {
			out = append(out, choice)
		}
	}
	return out
}

func eliminationKey(questionID string, playerID string) string {
	return questionID + "\x00" + playerID
}

func eliminatedChoicesByQuestionPlayer(match mongomodel.Match) map[string]map[string]mongomodel.MatchEliminatedChoice {
	out := make(map[string]map[string]mongomodel.MatchEliminatedChoice)
	for _, eliminated := range match.EliminatedChoices {
		if _, ok := out[eliminated.QuestionID]; !ok {
			out[eliminated.QuestionID] = make(map[string]mongomodel.MatchEliminatedChoice)
		}
		out[eliminated.QuestionID][eliminated.PlayerID] = eliminated
	}
	return out
}

func viewerEliminatedChoices(
	playerID string,
	viewerPlayerID string,
	currentEliminations map[string]mongomodel.MatchEliminatedChoice,
) ([]string, []string) {
	if playerID != viewerPlayerID {
		return nil, nil
	}
	eliminated, ok := currentEliminations[playerID]
	if !ok {
		return nil, nil
	}
	return cloneStrings(eliminated.Choices), cloneStrings(eliminated.SourceSitoneNames)
}

func choiceEliminated(match mongomodel.Match, questionID string, playerID string, choice string) bool {
	choice = strings.ToUpper(strings.TrimSpace(choice))
	for _, eliminated := range match.EliminatedChoices {
		if eliminated.QuestionID != questionID || eliminated.PlayerID != playerID {
			continue
		}
		for _, eliminatedChoice := range eliminated.Choices {
			if eliminatedChoice == choice {
				return true
			}
		}
	}
	return false
}
