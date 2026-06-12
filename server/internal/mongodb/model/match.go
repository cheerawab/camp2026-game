package model

import "time"

const MatchesCollection = "matches"

const (
	MatchStatusWaiting   = "waiting"
	MatchStatusActive    = "active"
	MatchStatusCompleted = "completed"
)

const (
	MatchPhaseAnswering = "answering"
	MatchPhaseRevealing = "revealing"
)

type Match struct {
	ID                   string        `bson:"_id"`
	Code                 string        `bson:"code"`
	Status               string        `bson:"status"`
	Phase                string        `bson:"phase,omitempty"`
	HostPlayerID         string        `bson:"host_player_id"`
	Players              []MatchPlayer `bson:"players"`
	QuestionIDs          []string      `bson:"question_ids,omitempty"`
	CurrentQuestionIndex int           `bson:"current_question_index,omitempty"`
	RoundStartedAt       time.Time     `bson:"round_started_at,omitempty"`
	RoundEndsAt          time.Time     `bson:"round_ends_at,omitempty"`
	RevealEndsAt         time.Time     `bson:"reveal_ends_at,omitempty"`
	CreatedAt            time.Time     `bson:"created_at"`
	StartedAt            time.Time     `bson:"started_at,omitempty"`
	CompletedAt          time.Time     `bson:"completed_at,omitempty"`
}

type MatchPlayer struct {
	PlayerID  string   `bson:"player_id"`
	Nickname  string   `bson:"nickname"`
	Ready     bool     `bson:"ready"`
	Score     int      `bson:"score"`
	SitoneIDs []string `bson:"sitone_ids,omitempty"`
}
