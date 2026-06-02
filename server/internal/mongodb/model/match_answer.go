package model

import "time"

const MatchAnswersCollection = "match_answers"

type MatchAnswer struct {
	ID            string    `bson:"_id"`
	MatchID       string    `bson:"match_id"`
	PlayerID      string    `bson:"player_id"`
	QuestionID    string    `bson:"question_id"`
	Choice        string    `bson:"choice"`
	Correct       bool      `bson:"correct"`
	Score         int       `bson:"score"`
	ElapsedMillis int64     `bson:"elapsed_ms"`
	AnsweredAt    time.Time `bson:"answered_at"`
}
