package apimodel

type QRCodeResponse struct {
	Token    string   `json:"token" example:"player_qr_token"`
	ImageURL string   `json:"imageUrl" example:"https://example.test/qrcode/player_qr_token.png"`
	Purposes []string `json:"purposes" example:"identity,match_pairing"`
}

type QRCodeScanRequest struct {
	Token   string `json:"token" validate:"required,min=8,max=200" example:"player_qr_token"`
	Context string `json:"context" validate:"required,oneof=staff_verification match_pairing" enums:"match_pairing" example:"match_pairing"`
}

type QRCodeScanResponse struct {
	Kind             string   `json:"kind" example:"player"`
	PlayerID         string   `json:"playerId" example:"player_01HR9Z7E2Z2VJ2QZ4P4Z"`
	DisplayName      string   `json:"displayName" example:"Alice"`
	AvailableActions []string `json:"availableActions" example:"create_match_pairing"`
}
