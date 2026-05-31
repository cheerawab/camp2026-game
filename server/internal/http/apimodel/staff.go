package apimodel

type StaffGrantRewardRequest struct {
	TargetQRCodeToken   string   `json:"targetQrCodeToken" validate:"required,min=8,max=200" example:"player_qr_token"`
	Reason              string   `json:"reason" validate:"required,max=120" example:"camp_activity_reward"`
	OpenPower           int      `json:"openPower" validate:"min=0,max=10000" example:"100"`
	SitoneDefinitionIDs []string `json:"sitoneDefinitionIds,omitempty" validate:"omitempty,dive,required" example:"sitone-engineering"`
	ItemDefinitionIDs   []string `json:"itemDefinitionIds,omitempty" validate:"omitempty,dive,required" example:"item-camp-sticker"`
}

type StaffGrantRewardResponse struct {
	GrantID  string `json:"grantId" example:"grant_01HR9Z7E2Z2VJ2QZ4P4Z"`
	PlayerID string `json:"playerId" example:"player_01HR9Z7E2Z2VJ2QZ4P4Z"`
	Reward   Reward `json:"reward"`
}

type StaffActivityVerificationRequest struct {
	TargetQRCodeToken string `json:"targetQrCodeToken" validate:"required,min=8,max=200" example:"player_qr_token"`
	ActivityCode      string `json:"activityCode" validate:"required,max=80" example:"booth-linux-101"`
	MissionID         string `json:"missionId,omitempty" example:"mission_activity_linux_101"`
}

type StaffActivityVerificationResponse struct {
	VerificationID string `json:"verificationId" example:"verify_01HR9Z7E2Z2VJ2QZ4P4Z"`
	PlayerID       string `json:"playerId" example:"player_01HR9Z7E2Z2VJ2QZ4P4Z"`
	MissionID      string `json:"missionId,omitempty" example:"mission_activity_linux_101"`
	Reward         Reward `json:"reward"`
}
