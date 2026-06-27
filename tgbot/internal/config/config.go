package config

import (
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"

	"github.com/sitcon-tw/camp2026-game/tgbot/internal/telegram"
)

const (
	defaultMongoURI          = "mongodb://camp2026:camp2026@localhost:27017/camp2026?authSource=admin"
	defaultMongoDatabase     = "camp2026"
	defaultLogLevel          = "info"
	defaultPollTimeout       = 50 * time.Second
	defaultRequestTTL        = 10 * time.Minute
	defaultHTTPClientTimeout = 60 * time.Second
	defaultShutdownTimeout   = 10 * time.Second
)

var defaultInitialSitoneIDs = []string{
	"stone_explorer_base",
	"stone_inspiration_base",
	"stone_resonance_base",
	"stone_engineering_base",
	"stone_entertainment_base",
}

type Config struct {
	BotToken          string
	APIBaseURL        string
	LoginBaseURL      string
	GroupTeamMap      map[int64]string
	InitialSitoneIDs  []string
	MongoURI          string
	MongoDatabase     string
	LogLevel          slog.Level
	PollTimeout       time.Duration
	RequestTTL        time.Duration
	HTTPClientTimeout time.Duration
	ShutdownTimeout   time.Duration
}

func Load() (Config, error) {
	loadDotEnv()

	cfg := Config{
		BotToken:          stringValue("TELEGRAM_BOT_TOKEN", ""),
		APIBaseURL:        stringValue("TELEGRAM_API_BASE_URL", telegram.DefaultAPIBaseURL),
		LoginBaseURL:      stringValue("APP_LOGIN_BASE_URL", ""),
		MongoURI:          stringValue("MONGODB_URI", defaultMongoURI),
		MongoDatabase:     stringValue("MONGODB_DATABASE", defaultMongoDatabase),
		InitialSitoneIDs:  listValue("INITIAL_SITONE_IDS", defaultInitialSitoneIDs),
		PollTimeout:       durationValue("TG_POLL_TIMEOUT", defaultPollTimeout),
		RequestTTL:        durationValue("TG_LOGIN_REQUEST_TTL", defaultRequestTTL),
		HTTPClientTimeout: durationValue("TG_HTTP_CLIENT_TIMEOUT", defaultHTTPClientTimeout),
		ShutdownTimeout:   durationValue("SHUTDOWN_TIMEOUT", defaultShutdownTimeout),
	}

	groupTeamMap, err := ParseGroupTeamMap(stringValue("TG_GROUP_TEAM_MAP", ""))
	if err != nil {
		return Config{}, err
	}
	cfg.GroupTeamMap = groupTeamMap

	level, err := parseLogLevel(stringValue("LOG_LEVEL", defaultLogLevel))
	if err != nil {
		return Config{}, err
	}
	cfg.LogLevel = level

	if err := cfg.validate(); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func loadDotEnv() {
	_ = godotenv.Load("tgbot/.env")
	_ = godotenv.Load(".env")
}

func ParseGroupTeamMap(value string) (map[int64]string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, errors.New("TG_GROUP_TEAM_MAP is required")
	}

	result := make(map[int64]string)
	seenTeams := make(map[string]int64)
	for _, rawEntry := range strings.Split(value, ",") {
		entry := strings.TrimSpace(rawEntry)
		if entry == "" {
			continue
		}
		chatIDRaw, teamID, ok := strings.Cut(entry, "=")
		if !ok {
			return nil, fmt.Errorf("invalid TG_GROUP_TEAM_MAP entry %q, expected chat_id=team_id", entry)
		}
		chatID, err := strconv.ParseInt(strings.TrimSpace(chatIDRaw), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid telegram chat id %q: %w", strings.TrimSpace(chatIDRaw), err)
		}
		teamID = strings.TrimSpace(teamID)
		if teamID == "" {
			return nil, fmt.Errorf("missing team id for telegram chat %d", chatID)
		}
		if _, exists := result[chatID]; exists {
			return nil, fmt.Errorf("duplicate telegram chat id %d", chatID)
		}
		if previousChatID, exists := seenTeams[teamID]; exists {
			return nil, fmt.Errorf("duplicate team id %q for telegram chats %d and %d", teamID, previousChatID, chatID)
		}
		result[chatID] = teamID
		seenTeams[teamID] = chatID
	}
	if len(result) == 0 {
		return nil, errors.New("TG_GROUP_TEAM_MAP must contain at least one chat_id=team_id entry")
	}
	return result, nil
}

func (cfg Config) validate() error {
	var errs []error
	if strings.TrimSpace(cfg.BotToken) == "" {
		errs = append(errs, errors.New("TELEGRAM_BOT_TOKEN is required"))
	}
	if err := validateURL("TELEGRAM_API_BASE_URL", cfg.APIBaseURL); err != nil {
		errs = append(errs, err)
	}
	if err := validateURL("APP_LOGIN_BASE_URL", cfg.LoginBaseURL); err != nil {
		errs = append(errs, err)
	}
	if strings.TrimSpace(cfg.MongoURI) == "" {
		errs = append(errs, errors.New("MONGODB_URI is required"))
	}
	if strings.TrimSpace(cfg.MongoDatabase) == "" {
		errs = append(errs, errors.New("MONGODB_DATABASE is required"))
	}
	if cfg.PollTimeout <= 0 {
		errs = append(errs, errors.New("TG_POLL_TIMEOUT must be positive"))
	}
	if cfg.RequestTTL <= 0 {
		errs = append(errs, errors.New("TG_LOGIN_REQUEST_TTL must be positive"))
	}
	if cfg.HTTPClientTimeout <= cfg.PollTimeout {
		errs = append(errs, errors.New("TG_HTTP_CLIENT_TIMEOUT must be greater than TG_POLL_TIMEOUT"))
	}
	if cfg.ShutdownTimeout <= 0 {
		errs = append(errs, errors.New("SHUTDOWN_TIMEOUT must be positive"))
	}
	if len(cfg.GroupTeamMap) == 0 {
		errs = append(errs, errors.New("TG_GROUP_TEAM_MAP is required"))
	}
	if len(cfg.InitialSitoneIDs) == 0 {
		errs = append(errs, errors.New("INITIAL_SITONE_IDS must include at least one sitone id"))
	}
	return errors.Join(errs...)
}

func validateURL(key string, value string) error {
	value = strings.TrimSpace(value)
	if value == "" {
		return fmt.Errorf("%s is required", key)
	}
	parsed, err := url.Parse(value)
	if err != nil {
		return fmt.Errorf("%s must be a valid URL: %w", key, err)
	}
	if parsed.Scheme == "" || parsed.Host == "" {
		return fmt.Errorf("%s must include scheme and host", key)
	}
	return nil
}

func stringValue(key string, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}

func listValue(key string, fallback []string) []string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return append([]string(nil), fallback...)
	}
	parts := strings.Split(value, ",")
	out := make([]string, 0, len(parts))
	seen := make(map[string]struct{}, len(parts))
	for _, part := range parts {
		item := strings.TrimSpace(part)
		if item == "" {
			continue
		}
		if _, ok := seen[item]; ok {
			continue
		}
		seen[item] = struct{}{}
		out = append(out, item)
	}
	return out
}

func durationValue(key string, fallback time.Duration) time.Duration {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	parsed, err := time.ParseDuration(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func parseLogLevel(value string) (slog.Level, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "debug":
		return slog.LevelDebug, nil
	case "info", "":
		return slog.LevelInfo, nil
	case "warn", "warning":
		return slog.LevelWarn, nil
	case "error":
		return slog.LevelError, nil
	default:
		return slog.LevelInfo, fmt.Errorf("LOG_LEVEL must be one of debug, info, warn, error")
	}
}
