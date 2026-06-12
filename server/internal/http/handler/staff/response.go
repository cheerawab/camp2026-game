package staff

type CreateRewardRequest struct {
	QRCodeToken string `json:"qrcodeToken" validate:"required,min=4,max=512" example:"qr_token_123456"`
	Kind        string `json:"kind" validate:"required,oneof=item sitone" example:"sitone"`
	RefID       string `json:"refId" validate:"required,min=1,max=128" example:"stone_engineering_base"`
	Quantity    int    `json:"quantity" validate:"required,min=1,max=99" example:"1"`
}

type CreateRewardResponse struct {
	RewardID string               `json:"rewardId" example:"staff_reward_01HXK2P9ATJ5S2YV8C2J4Q0M"`
	Player   RewardPlayerResponse `json:"player"`
	Reward   RewardResponse       `json:"reward"`
}

type RewardPlayerResponse struct {
	PlayerID string             `json:"playerId" example:"7H9K2Q"`
	Nickname string             `json:"nickname" example:"Alice"`
	Team     RewardTeamResponse `json:"team"`
}

type RewardTeamResponse struct {
	TeamID string `json:"teamId" example:"8M4RXP"`
	Name   string `json:"name" example:"Blue Team"`
}

type RewardResponse struct {
	Kind     string `json:"kind" example:"sitone"`
	ID       string `json:"id" example:"stone_engineering_base"`
	Name     string `json:"name" example:"工程型小石"`
	Quantity int    `json:"quantity" example:"1"`
}
