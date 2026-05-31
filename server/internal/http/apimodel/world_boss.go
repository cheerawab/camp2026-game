package apimodel

type WorldBossListResponse struct {
	Bosses []WorldBossSummary `json:"bosses"`
}

type WorldBossDetailResponse struct {
	Boss    WorldBossSummary  `json:"boss"`
	Rewards []WorldBossReward `json:"rewards"`
}

type WorldBossSummary struct {
	BossID            string `json:"bossId" example:"boss_layer_1"`
	Name              string `json:"name" example:"Knowledge Core"`
	Layer             int    `json:"layer" example:"1"`
	HP                int    `json:"hp" example:"4500"`
	MaxHP             int    `json:"maxHp" example:"10000"`
	RemainingAttempts int    `json:"remainingAttempts" example:"2"`
	Status            string `json:"status" example:"active"`
}

type WorldBossReward struct {
	RewardID string `json:"rewardId" example:"wb_reward_layer_1"`
	Status   string `json:"status" example:"locked"`
	Reward   Reward `json:"reward"`
}

type WorldBossMatchCreateRequest struct {
	SitoneIDs []string `json:"sitoneIds" validate:"max=5,dive,required" example:"sitone_01HR9Z7E2Z2VJ2QZ4P4Z"`
}

type WorldBossMatchResponse struct {
	MatchID           string        `json:"matchId" example:"match_01HR9Z7E2Z2VJ2QZ4P4Z"`
	BossID            string        `json:"bossId" example:"boss_layer_1"`
	RemainingAttempts int           `json:"remainingAttempts" example:"1"`
	Question          MatchQuestion `json:"question"`
}
