package me

type TeamResponse struct {
	TeamID string `json:"teamId" example:"8M4RXP"`
	Name   string `json:"name" example:"Blue Team"`
}

type StatusResponse struct {
	PlayerID  string       `json:"playerId" example:"7H9K2Q"`
	Nickname  string       `json:"nickname" example:"Alice"`
	Team      TeamResponse `json:"team"`
	OpenPower int          `json:"openPower" example:"1280"`
	AvatarURL string       `json:"avatarUrl,omitempty" example:"https://example.test/avatar/alice.png"`
}

type HomeResponse struct {
	Player   StatusResponse       `json:"player"`
	Summary  HomeSummaryResponse  `json:"summary"`
	TeamRank *TeamRankResponse    `json:"teamRank,omitempty"`
	Actions  []HomeActionResponse `json:"actions"`
}

type HomeSummaryResponse struct {
	OpenPower   int `json:"openPower" example:"1280"`
	SitoneCount int `json:"sitoneCount" example:"8"`
	ItemCount   int `json:"itemCount" example:"29"`
}

type TeamRankResponse struct {
	Type          string `json:"type" example:"open_power"`
	Rank          int    `json:"rank" example:"2"`
	TeamID        string `json:"teamId" example:"8M4RXP"`
	Name          string `json:"name" example:"Blue Team"`
	Score         int    `json:"score" example:"1188"`
	GapToPrevious int    `json:"gapToPrevious" example:"72"`
}

type HomeActionResponse struct {
	ID      string `json:"id" example:"battle"`
	Label   string `json:"label" example:"知識王戰"`
	Enabled bool   `json:"enabled" example:"true"`
}

type QRCodeResponse struct {
	QRCodeToken string `json:"qrcodeToken" example:"qr_token_123456"`
}

type SitoneListResponse struct {
	Sitones []PlayerSitoneResponse `json:"sitones"`
}

type PlayerSitoneResponse struct {
	ID       string         `json:"id" example:"owned-sitone-001"`
	SitoneID string         `json:"sitoneId" example:"sitone-engineering"`
	Quantity int            `json:"quantity" example:"1"`
	Sitone   SitoneResponse `json:"sitone"`
}

type SitoneResponse struct {
	ID          string `json:"id" example:"sitone-engineering"`
	Name        string `json:"name" example:"工程型小石"`
	Type        string `json:"type" example:"engineering"`
	Rarity      string `json:"rarity" example:"base"`
	Style       string `json:"style" example:"default"`
	Description string `json:"description" example:"修 bug、分享解法、完成技術任務。"`
}

type ItemListResponse struct {
	Items []PlayerItemResponse `json:"items"`
}

type PlayerItemResponse struct {
	ID       string       `json:"id" example:"owned-item-001"`
	ItemID   string       `json:"itemId" example:"item-crafting-fragment"`
	Quantity int          `json:"quantity" example:"3"`
	Item     ItemResponse `json:"item"`
}

type ItemResponse struct {
	ID          string `json:"id" example:"item-crafting-fragment"`
	Name        string `json:"name" example:"合成碎片"`
	Type        string `json:"type" example:"material"`
	Rarity      string `json:"rarity" example:"common"`
	Description string `json:"description" example:"小石造型合成使用的基礎素材。"`
}
