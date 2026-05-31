package apimodel

type MissionListResponse struct {
	Tabs     []MissionTab     `json:"tabs"`
	Missions []MissionSummary `json:"missions"`
}

type BingoBoardListResponse struct {
	Boards []BingoBoardSummary `json:"boards"`
}

type BingoBoardSummary struct {
	BoardID  string        `json:"boardId" example:"board_day_1"`
	Title    string        `json:"title" example:"Day 1 Bingo"`
	Category string        `json:"category" example:"daily"`
	Status   string        `json:"status" example:"active"`
	Cells    []BingoCell   `json:"cells"`
	Rewards  []BingoReward `json:"rewards"`
}

type BingoCell struct {
	Row     int            `json:"row" example:"0"`
	Column  int            `json:"column" example:"0"`
	Mission MissionSummary `json:"mission"`
}

type BingoReward struct {
	LineRewardID string `json:"lineRewardId" example:"line_reward_row_0"`
	LineKey      string `json:"lineKey" example:"row_0"`
	Status       string `json:"status" example:"claimable"`
	Reward       Reward `json:"reward"`
}

type MissionTab struct {
	Key   string `json:"key" example:"daily"`
	Label string `json:"label" example:"Daily"`
	Count int    `json:"count" example:"4"`
}

type MissionSummary struct {
	ID          string          `json:"id" example:"mission_daily_match_3"`
	Tab         string          `json:"tab" example:"daily"`
	Title       string          `json:"title" example:"Answer three match questions"`
	Description string          `json:"description" example:"Complete three Knowledge King questions today."`
	Status      string          `json:"status" example:"claimable"`
	Progress    MissionProgress `json:"progress"`
	Rewards     Reward          `json:"rewards"`
}

type MissionProgress struct {
	Current int `json:"current" example:"3"`
	Target  int `json:"target" example:"3"`
}

type MissionCompleteRequest struct {
	Flag              string `json:"flag,omitempty" validate:"omitempty,min=3,max=80" example:"CAMP2026-HELLO"`
	StaffQRCodeToken  string `json:"staffQrCodeToken,omitempty" validate:"omitempty,min=8,max=200" swaggerignore:"true"`
	ClientCompletedAt string `json:"clientCompletedAt,omitempty" example:"2026-07-24T10:30:00+08:00"`
}

type MissionCompleteResponse struct {
	Mission MissionSummary `json:"mission"`
	Reward  Reward         `json:"reward"`
}

type BingoLineRewardClaimResponse struct {
	LineRewardID string `json:"lineRewardId" example:"line_reward_row_0"`
	Status       string `json:"status" example:"claimed"`
	Reward       Reward `json:"reward"`
}
