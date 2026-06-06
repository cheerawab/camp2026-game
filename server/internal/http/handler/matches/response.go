package matches

import "time"

type CreateMatchResponse = MatchStateResponse
type JoinMatchResponse = MatchStateResponse
type ReadyMatchResponse = MatchStateResponse

type JoinMatchRequest struct {
	Code string `json:"code" validate:"required,min=4,max=16" example:"ABC123"`
}

type JoinByQRRequest struct {
	QRCodeToken string `json:"qrcodeToken" validate:"required,min=4,max=512" example:"qr_token_123456"`
}

type AnswerRequest struct {
	QuestionID string `json:"questionId" validate:"required" example:"quiz-001"`
	Choice     string `json:"choice" validate:"required,oneof=A B C D" example:"A"`
}

type AnswerAcceptedResponse struct {
	Accepted bool `json:"accepted" example:"true"`
}

type MatchStateResponse struct {
	MatchID              string                 `json:"matchId" example:"match_7H9K2Q"`
	Code                 string                 `json:"code,omitempty" example:"ABC123"`
	Status               string                 `json:"status" example:"active"`
	HostPlayerID         string                 `json:"hostPlayerId" example:"7H9K2Q"`
	Players              []MatchPlayerResponse  `json:"players"`
	CurrentQuestionIndex int                    `json:"currentQuestionIndex,omitempty" example:"0"`
	QuestionCount        int                    `json:"questionCount,omitempty" example:"10"`
	CurrentQuestion      *MatchQuestionResponse `json:"currentQuestion,omitempty"`
	RoundStartedAt       *time.Time             `json:"roundStartedAt,omitempty"`
	RoundEndsAt          *time.Time             `json:"roundEndsAt,omitempty"`
	CreatedAt            time.Time              `json:"createdAt"`
	StartedAt            *time.Time             `json:"startedAt,omitempty"`
	CompletedAt          *time.Time             `json:"completedAt,omitempty"`
	Results              []MatchQuestionResult  `json:"results,omitempty"`
}

type MatchPlayerResponse struct {
	PlayerID                string `json:"playerId" example:"7H9K2Q"`
	Nickname                string `json:"nickname" example:"Alice"`
	Ready                   bool   `json:"ready" example:"true"`
	AnsweredCurrentQuestion bool   `json:"answeredCurrentQuestion,omitempty" example:"true"`
	Score                   *int   `json:"score,omitempty" example:"850"`
	OpenPowerReward         *int   `json:"openPowerReward,omitempty" example:"105"`
}

type MatchQuestionResponse struct {
	QuestionID string `json:"questionId" example:"quiz-001"`
	Prompt     string `json:"prompt" example:"Which command initializes a new Git repository?"`
	ChoiceA    string `json:"choiceA" example:"git init"`
	ChoiceB    string `json:"choiceB" example:"git clone"`
	ChoiceC    string `json:"choiceC" example:"git status"`
	ChoiceD    string `json:"choiceD" example:"git add"`
}

type MatchQuestionResult struct {
	QuestionID    string                `json:"questionId" example:"quiz-001"`
	Prompt        string                `json:"prompt" example:"Which command initializes a new Git repository?"`
	ChoiceA       string                `json:"choiceA" example:"git init"`
	ChoiceB       string                `json:"choiceB" example:"git clone"`
	ChoiceC       string                `json:"choiceC" example:"git status"`
	ChoiceD       string                `json:"choiceD" example:"git add"`
	CorrectChoice string                `json:"correctChoice" example:"A"`
	Explanation   string                `json:"explanation" example:"git init creates a new repository."`
	Answers       []MatchAnswerResponse `json:"answers"`
}

type MatchAnswerResponse struct {
	PlayerID      string     `json:"playerId" example:"7H9K2Q"`
	Choice        string     `json:"choice,omitempty" example:"A"`
	Correct       bool       `json:"correct" example:"true"`
	Score         int        `json:"score" example:"150"`
	ElapsedMillis int64      `json:"elapsedMillis" example:"3200"`
	AnsweredAt    *time.Time `json:"answeredAt,omitempty"`
}
