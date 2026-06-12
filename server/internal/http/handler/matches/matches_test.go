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

func TestMatchOpenPowerRewardRequiresClearWinner(t *testing.T) {
	match := mongomodel.Match{
		Players: []mongomodel.MatchPlayer{
			{PlayerID: "P1", Score: 340},
			{PlayerID: "P2", Score: 340},
		},
	}

	if matchHasClearWinner(match) {
		t.Fatal("expected tied match to have no clear winner")
	}
	for _, player := range match.Players {
		if got := matchOpenPowerReward(match, player); got != 0 {
			t.Fatalf("expected tied match reward 0 for %s, got %d", player.PlayerID, got)
		}
	}

	match.Players[1].Score = 330
	if !matchHasClearWinner(match) {
		t.Fatal("expected match with different scores to have a clear winner")
	}
	if got := matchOpenPowerReward(match, match.Players[0]); got != 54 {
		t.Fatalf("expected winner reward 54, got %d", got)
	}
	if got := matchOpenPowerReward(match, match.Players[1]); got != 53 {
		t.Fatalf("expected non-tied opponent reward 53, got %d", got)
	}
}

func TestMatchRewardKeysUseMatchAndPlayer(t *testing.T) {
	if got := matchRewardRecordID("match_123", "P1"); got != "open_power_reward_match_123_P1" {
		t.Fatalf("unexpected reward record id: %q", got)
	}
	if got := matchRewardSource("match_123", "P1"); got != "quiz_match:match_123:player:P1" {
		t.Fatalf("unexpected reward source: %q", got)
	}
}

func TestAllPlayersReady(t *testing.T) {
	match := mongomodel.Match{
		Players: []mongomodel.MatchPlayer{
			{PlayerID: "P1", Ready: true, SitoneIDs: []string{"stone_engineering_base"}},
			{PlayerID: "P2", Ready: true, SitoneIDs: []string{"stone_resonance_base"}},
		},
	}
	if !allPlayersReady(match) {
		t.Fatal("expected all players ready")
	}

	match.Players[1].Ready = false
	if allPlayersReady(match) {
		t.Fatal("expected not all players ready")
	}

	match.Players[1].Ready = true
	match.Players[1].SitoneIDs = nil
	if allPlayersReady(match) {
		t.Fatal("expected not all players ready without loadout")
	}
}

func TestNormalizeSitoneLoadoutAllowsDuplicateSlots(t *testing.T) {
	got, err := normalizeSitoneLoadout([]string{
		" stone_engineering_base ",
		"stone_engineering_base",
		"",
	})
	if err != nil {
		t.Fatalf("normalize sitone loadout: %v", err)
	}

	want := []string{"stone_engineering_base", "stone_engineering_base"}
	if len(got) != len(want) {
		t.Fatalf("expected %d sitones, got %#v", len(want), got)
	}
	for index := range want {
		if got[index] != want[index] {
			t.Fatalf("unexpected sitone at index %d: got %q want %q", index, got[index], want[index])
		}
	}
}

func TestMaxScoreThroughCurrentQuestion(t *testing.T) {
	match := mongomodel.Match{
		Status:               mongomodel.MatchStatusActive,
		QuestionIDs:          []string{"quiz-001", "quiz-002", "quiz-003"},
		CurrentQuestionIndex: 1,
	}

	if got := maxScorePerQuestion(); got != 175 {
		t.Fatalf("expected max score per question 175, got %d", got)
	}
	if got := maxScoreThroughCurrentQuestion(match); got != 350 {
		t.Fatalf("expected active max score 350, got %d", got)
	}

	match.Status = mongomodel.MatchStatusCompleted
	if got := maxScoreThroughCurrentQuestion(match); got != 525 {
		t.Fatalf("expected completed max score 525, got %d", got)
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
