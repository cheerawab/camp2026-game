package domain

import (
	"errors"
	"time"
)

const (
	PlayersCollection       = "players"
	PlayerSitonesCollection = "player_sitones"
	LoginRequestsCollection = "telegram_login_requests"
)

var ErrLoginRequestNotFound = errors.New("telegram login request not found")

type LoginRequest struct {
	ID             string    `bson:"_id"`
	TelegramUserID int64     `bson:"telegram_user_id"`
	TelegramChatID int64     `bson:"telegram_chat_id"`
	TeamID         string    `bson:"team_id"`
	CreatedAt      time.Time `bson:"created_at"`
	ExpiresAt      time.Time `bson:"expires_at"`
}

type Player struct {
	ID               string    `bson:"_id"`
	AuthToken        string    `bson:"auth_token"`
	QRCodeToken      string    `bson:"qrcode_token,omitempty"`
	Nickname         string    `bson:"nickname"`
	TeamID           string    `bson:"team_id,omitempty"`
	AvatarURL        string    `bson:"avatar_url,omitempty"`
	Role             string    `bson:"role,omitempty"`
	DefaultSitoneIDs []string  `bson:"default_sitone_ids,omitempty"`
	TelegramUserID   int64     `bson:"telegram_user_id,omitempty"`
	TelegramUsername string    `bson:"telegram_username,omitempty"`
	TelegramChatID   int64     `bson:"telegram_chat_id,omitempty"`
	CreatedAt        time.Time `bson:"created_at,omitempty"`
	UpdatedAt        time.Time `bson:"updated_at,omitempty"`
}

type CreatePlayerInput struct {
	PlayerID         string
	AuthToken        string
	QRCodeToken      string
	Nickname         string
	TeamID           string
	InitialSitoneIDs []string
	TelegramUserID   int64
	TelegramUsername string
	TelegramChatID   int64
	Now              time.Time
}

type PlayerSitone struct {
	ID       string `bson:"_id"`
	PlayerID string `bson:"player_id"`
	SitoneID string `bson:"sitone_id"`
	Quantity int    `bson:"quantity"`
}
