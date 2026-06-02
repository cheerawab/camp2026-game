package matches

import (
	"testing"
	"time"

	"github.com/sitcon-tw/camp2026-game/internal/content"
	mongomodel "github.com/sitcon-tw/camp2026-game/internal/mongodb/model"
)

func TestScoreAnswer(t *testing.T) {
	question := content.QuizQuestion{CorrectChoice: "A"}
	answeredAt := time.Date(2026, 6, 2, 10, 0, 3, 0, time.UTC)
	roundEndsAt := time.Date(2026, 6, 2, 10, 0, 15, 0, time.UTC)

	correct, score := scoreAnswer(question, "A", answeredAt, roundEndsAt)
	if !correct {
		t.Fatal("expected answer to be correct")
	}
	if score != 160 {
		t.Fatalf("expected score 160, got %d", score)
	}

	correct, score = scoreAnswer(question, "B", answeredAt, roundEndsAt)
	if correct {
		t.Fatal("expected answer to be incorrect")
	}
	if score != 0 {
		t.Fatalf("expected incorrect score 0, got %d", score)
	}
}

func TestOpenPowerReward(t *testing.T) {
	if got := openPowerReward(850); got != 105 {
		t.Fatalf("expected open power reward 105, got %d", got)
	}
}

func TestAllPlayersReady(t *testing.T) {
	match := mongomodel.Match{
		Players: []mongomodel.MatchPlayer{
			{PlayerID: "P1", Ready: true},
			{PlayerID: "P2", Ready: true},
		},
	}
	if !allPlayersReady(match) {
		t.Fatal("expected all players ready")
	}

	match.Players[1].Ready = false
	if allPlayersReady(match) {
		t.Fatal("expected not all players ready")
	}
}

func TestPickQuestionIDs(t *testing.T) {
	handler := New(Dependencies{Content: loadTestContent(t)})

	ids, err := handler.pickQuestionIDs()
	if err != nil {
		t.Fatalf("pick question ids: %v", err)
	}
	if len(ids) != matchQuestionCount {
		t.Fatalf("expected %d question ids, got %d", matchQuestionCount, len(ids))
	}

	seen := make(map[string]struct{}, len(ids))
	for _, id := range ids {
		if _, ok := seen[id]; ok {
			t.Fatalf("expected unique question ids, got duplicate %q", id)
		}
		seen[id] = struct{}{}
		if _, ok := handler.content.GetQuizQuestion(id); !ok {
			t.Fatalf("expected question %q to exist", id)
		}
	}
}

func TestMatchResultsRevealCompletedAnswers(t *testing.T) {
	handler := New(Dependencies{Content: loadTestContent(t)})
	match := mongomodel.Match{
		ID:          "match_123",
		Status:      mongomodel.MatchStatusCompleted,
		QuestionIDs: []string{"quiz-001"},
		Players: []mongomodel.MatchPlayer{
			{PlayerID: "P1", Nickname: "Alice", Score: 160},
			{PlayerID: "P2", Nickname: "Bob", Score: 0},
		},
	}
	answeredAt := time.Date(2026, 6, 2, 10, 0, 3, 0, time.UTC)
	answers := map[string]map[string]mongomodel.MatchAnswer{
		"quiz-001": {
			"P1": {
				PlayerID:      "P1",
				QuestionID:    "quiz-001",
				Choice:        "A",
				Correct:       true,
				Score:         160,
				ElapsedMillis: 3000,
				AnsweredAt:    answeredAt,
			},
		},
	}

	results, err := handler.matchResults(match, answers)
	if err != nil {
		t.Fatalf("match results: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].CorrectChoice == "" || results[0].Explanation == "" {
		t.Fatalf("expected completed result to reveal correct choice and explanation: %#v", results[0])
	}
	if len(results[0].Answers) != 2 {
		t.Fatalf("expected 2 player answer rows, got %d", len(results[0].Answers))
	}
	if results[0].Answers[0].Choice != "A" || !results[0].Answers[0].Correct {
		t.Fatalf("expected first answer to be revealed, got %#v", results[0].Answers[0])
	}
}

func TestBrokerPublishesEvents(t *testing.T) {
	broker := NewBroker()
	events, unsubscribe := broker.Subscribe("match_123")
	defer unsubscribe()

	broker.Publish("match_123", Event{Name: "player_answered"})

	select {
	case event := <-events:
		if event.Name != "player_answered" {
			t.Fatalf("expected player_answered event, got %q", event.Name)
		}
	case <-time.After(time.Second):
		t.Fatal("expected event")
	}
}

func loadTestContent(t *testing.T) *content.Store {
	t.Helper()

	store, err := content.Load("../../../../content")
	if err != nil {
		t.Fatalf("load test content: %v", err)
	}
	return store
}
