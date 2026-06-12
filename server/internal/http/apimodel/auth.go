package apimodel

type AuthLoginRequest struct {
	Token string `json:"token" validate:"required,max=512" example:"token-1"`
}

type AuthLoginResponse struct {
	Player AuthPlayerSummary `json:"player"`
}

type AuthPlayerSummary struct {
	PlayerID  string          `json:"playerId" example:"7H9K2Q"`
	Nickname  string          `json:"nickname" example:"Alice"`
	Team      AuthTeamSummary `json:"team"`
	OpenPower int             `json:"openPower" example:"1280"`
	AvatarURL string          `json:"avatarUrl,omitempty" example:"https://example.test/avatar/alice.png"`
	Role      string          `json:"role,omitempty" example:"staff"`
}

type AuthTeamSummary struct {
	TeamID string `json:"teamId" example:"8M4RXP"`
	Name   string `json:"name" example:"Blue Team"`
}
