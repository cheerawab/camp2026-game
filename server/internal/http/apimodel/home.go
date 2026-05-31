package apimodel

type UserStateResponse struct {
	Player   UserStatePlayer   `json:"player"`
	Stats    UserStateStats    `json:"stats"`
	Features UserFeatureStates `json:"features"`
	Loadout  []SitoneSummary   `json:"loadout"`
}

type UserStatePlayer struct {
	PlayerID  string `json:"playerId" example:"player_01HR9Z7E2Z2VJ2QZ4P4Z"`
	TeamID    string `json:"teamId" example:"team_blue"`
	Nickname  string `json:"nickname" example:"Alice"`
	OpenPower int    `json:"openPower" example:"1280"`
}

type UserStateStats struct {
	OpenPower                   int `json:"openPower" example:"1280"`
	CompletedMissionCount       int `json:"completedMissionCount" example:"7"`
	ClaimableMissionRewardCount int `json:"claimableMissionRewardCount" example:"2"`
	MatchCount                  int `json:"matchCount" example:"12"`
	MatchWinCount               int `json:"matchWinCount" example:"8"`
	OwnedSitoneCount            int `json:"ownedSitoneCount" example:"5"`
	OwnedItemCount              int `json:"ownedItemCount" example:"3"`
}

type UserFeatureStates struct {
	Bingo     UserBingoState     `json:"bingo"`
	Matches   UserMatchState     `json:"matches"`
	WorldBoss UserWorldBossState `json:"worldBoss"`
	Storage   UserStorageState   `json:"storage"`
	QRCode    UserQRCodeState    `json:"qrCode"`
}

type UserBingoState struct {
	Enabled               bool `json:"enabled" example:"true"`
	ActiveMissionCount    int  `json:"activeMissionCount" example:"4"`
	CompletedMissionCount int  `json:"completedMissionCount" example:"7"`
	ClaimableRewardCount  int  `json:"claimableRewardCount" example:"2"`
}

type UserMatchState struct {
	Enabled             bool `json:"enabled" example:"true"`
	PendingPairingCount int  `json:"pendingPairingCount" example:"1"`
	ActiveMatchCount    int  `json:"activeMatchCount" example:"0"`
	CompletedMatchCount int  `json:"completedMatchCount" example:"12"`
	WinCount            int  `json:"winCount" example:"8"`
}

type UserWorldBossState struct {
	Enabled               bool   `json:"enabled" example:"true"`
	ActiveBossID          string `json:"activeBossId,omitempty" example:"boss_layer_1"`
	RemainingAttemptCount int    `json:"remainingAttemptCount" example:"3"`
	ClaimableRewardCount  int    `json:"claimableRewardCount" example:"0"`
}

type UserStorageState struct {
	Enabled              bool `json:"enabled" example:"true"`
	SitoneCount          int  `json:"sitoneCount" example:"5"`
	ItemCount            int  `json:"itemCount" example:"3"`
	CraftableRecipeCount int  `json:"craftableRecipeCount" example:"1"`
}

type UserQRCodeState struct {
	Enabled        bool `json:"enabled" example:"true"`
	HasActiveToken bool `json:"hasActiveToken" example:"true"`
}
