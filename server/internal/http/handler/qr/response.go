package qr

type ResolveRequest struct {
	QRCodeToken string `json:"qrcodeToken" validate:"required,min=4,max=512" example:"qr_token_123456"`
}

type ResolveResponse struct {
	Player PlayerSummaryResponse `json:"player"`
}

type PlayerSummaryResponse struct {
	PlayerID  string       `json:"playerId" example:"7H9K2Q"`
	Nickname  string       `json:"nickname" example:"Alice"`
	Team      TeamResponse `json:"team"`
	OpenPower int          `json:"openPower" example:"1280"`
	AvatarURL string       `json:"avatarUrl,omitempty" example:"https://example.test/avatar/alice.png"`
}

type TeamResponse struct {
	TeamID string `json:"teamId" example:"8M4RXP"`
	Name   string `json:"name" example:"Blue Team"`
}
