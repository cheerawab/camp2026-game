package apimodel

type MatchPairingRequest struct {
	TargetQRCodeToken string `json:"targetQrCodeToken" validate:"required,min=8,max=200" example:"player_qr_token"`
}

type MatchPairingResponse struct {
	PairingID string `json:"pairingId" example:"pair_01HR9Z7E2Z2VJ2QZ4P4Z"`
	Status    string `json:"status" example:"matched"`
}

type MatchCreateRequest struct {
	Mode             string   `json:"mode" validate:"required,oneof=qr_duel offline_duel world_boss" example:"qr_duel"`
	PairingID        string   `json:"pairingId,omitempty" example:"pair_01HR9Z7E2Z2VJ2QZ4P4Z"`
	OpponentPlayerID string   `json:"opponentPlayerId,omitempty" example:"player_01HR9Z7E2Z2VJ2QZ4P4Z"`
	WorldBossID      string   `json:"worldBossId,omitempty" example:"boss_layer_1"`
	SitoneIDs        []string `json:"sitoneIds" validate:"max=5,dive,required" example:"sitone_01HR9Z7E2Z2VJ2QZ4P4Z"`
}

type MatchResponse struct {
	MatchID  string           `json:"matchId" example:"match_01HR9Z7E2Z2VJ2QZ4P4Z"`
	Mode     string           `json:"mode" example:"qr_duel"`
	Status   string           `json:"status" example:"answering"`
	Player   MatchParticipant `json:"player"`
	Opponent MatchParticipant `json:"opponent"`
	Question MatchQuestion    `json:"question"`
}

type MatchListResponse struct {
	Matches []MatchSummary `json:"matches"`
}

type MatchSummary struct {
	MatchID         string `json:"matchId" example:"match_01HR9Z7E2Z2VJ2QZ4P4Z"`
	Mode            string `json:"mode" example:"qr_duel"`
	Status          string `json:"status" example:"completed"`
	OpponentName    string `json:"opponentName" example:"Bob"`
	PlayerScore     int    `json:"playerScore" example:"320"`
	OpponentScore   int    `json:"opponentScore" example:"280"`
	CompletedAt     string `json:"completedAt,omitempty" example:"2026-07-24T10:35:00+08:00"`
	WorldBossID     string `json:"worldBossId,omitempty" example:"boss_layer_1"`
	DamageDealt     int    `json:"damageDealt,omitempty" example:"320"`
	OpenPowerGained int    `json:"openPowerGained" example:"80"`
}

type MatchParticipant struct {
	PlayerID    string `json:"playerId" example:"player_01HR9Z7E2Z2VJ2QZ4P4Z"`
	DisplayName string `json:"displayName" example:"Alice"`
	HP          int    `json:"hp" example:"100"`
	OpenPower   int    `json:"openPower" example:"1280"`
}

type MatchQuestion struct {
	QuestionID string        `json:"questionId" example:"question_001"`
	Prompt     string        `json:"prompt" example:"Which license is commonly used for open source projects?"`
	Choices    []MatchChoice `json:"choices"`
	TimeLimit  int           `json:"timeLimitSeconds" example:"30"`
}

type MatchChoice struct {
	ChoiceID string `json:"choiceId" example:"A"`
	Text     string `json:"text" example:"MIT License"`
}

type MatchAnswerSubmitRequest struct {
	QuestionID       string `json:"questionId" validate:"required" example:"question_001"`
	ChoiceID         string `json:"choiceId" validate:"required" example:"A"`
	ClientAnsweredAt string `json:"clientAnsweredAt,omitempty" example:"2026-07-24T10:31:00+08:00"`
}

type MatchAnswerSubmitResponse struct {
	Correct         bool             `json:"correct" example:"true"`
	CorrectChoiceID string           `json:"correctChoiceId" example:"A"`
	Explanation     string           `json:"explanation" example:"The MIT License is a permissive open source license."`
	OpenPowerGained int              `json:"openPowerGained" example:"80"`
	Battle          MatchBattleState `json:"battle"`
	NextQuestion    *MatchQuestion   `json:"nextQuestion,omitempty"`
}

type MatchBattleState struct {
	PlayerHP   int `json:"playerHp" example:"90"`
	OpponentHP int `json:"opponentHp" example:"65"`
	Round      int `json:"round" example:"2"`
}

type MatchWebSocketInfoResponse struct {
	MatchID string   `json:"matchId" example:"match_01HR9Z7E2Z2VJ2QZ4P4Z"`
	Events  []string `json:"events" example:"match.snapshot,answer.submit,answer.result,match.completed,error"`
}
