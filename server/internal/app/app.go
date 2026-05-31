package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/sitcon-tw/camp2026-game/internal/config"
	httpserver "github.com/sitcon-tw/camp2026-game/internal/http"
	"github.com/sitcon-tw/camp2026-game/internal/postgres"
)

type Application struct {
	Config     config.Config
	Log        *slog.Logger
	HTTPServer *http.Server
	DBPool     *pgxpool.Pool
}

func New(ctx context.Context) (*Application, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: cfg.LogLevel,
	}))

	dbPool, err := postgres.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}

	handler := httpserver.NewRouter(httpserver.Dependencies{
		Log:            log,
		RequestTimeout: cfg.HTTP.RequestTimeout,
		ReadinessCheck: postgres.ReadinessCheck(dbPool),
	})

	return &Application{
		Config:     cfg,
		Log:        log,
		HTTPServer: httpserver.NewServer(cfg.HTTP, handler),
		DBPool:     dbPool,
	}, nil
}
