package leaderboards

type ListResponse struct {
	Type          string             `json:"type" example:"open_power"`
	Teams         []TeamRankResponse `json:"teams"`
	CurrentTeam   *TeamRankResponse  `json:"currentTeam,omitempty"`
	GapToPrevious int                `json:"gapToPrevious" example:"72"`
}

type TeamRankResponse struct {
	Rank    int    `json:"rank" example:"2"`
	TeamID  string `json:"teamId" example:"8M4RXP"`
	Name    string `json:"name" example:"Blue Team"`
	Score   int    `json:"score" example:"1188"`
	Metric  string `json:"metric" example:"OP"`
	Current bool   `json:"current" example:"true"`
}
