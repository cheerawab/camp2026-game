package bot

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/sitcon-tw/camp2026-game/tgbot/internal/telegram"
)

type UpdateSource interface {
	GetUpdates(ctx context.Context, offset int, timeout time.Duration) ([]telegram.Update, error)
}

type Runner struct {
	source      UpdateSource
	service     *Service
	pollTimeout time.Duration
	log         *slog.Logger
}

func NewRunner(source UpdateSource, service *Service, pollTimeout time.Duration, log *slog.Logger) *Runner {
	if pollTimeout <= 0 {
		pollTimeout = 50 * time.Second
	}
	if log == nil {
		log = slog.Default()
	}
	return &Runner{
		source:      source,
		service:     service,
		pollTimeout: pollTimeout,
		log:         log,
	}
}

func (r *Runner) Run(ctx context.Context) error {
	offset := 0
	backoff := time.Second
	for ctx.Err() == nil {
		updates, err := r.source.GetUpdates(ctx, offset, r.pollTimeout)
		if err != nil {
			if ctx.Err() != nil && (errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded)) {
				return ctx.Err()
			}
			r.log.Warn("telegram getUpdates failed", "error", err)
			select {
			case <-time.After(backoff):
			case <-ctx.Done():
				return ctx.Err()
			}
			if backoff < 30*time.Second {
				backoff *= 2
			}
			continue
		}
		backoff = time.Second

		for _, update := range updates {
			if update.UpdateID >= offset {
				offset = update.UpdateID + 1
			}
			if err := r.service.HandleUpdate(ctx, update); err != nil {
				r.log.Error("telegram update handling failed",
					"update_id", update.UpdateID,
					"error", err,
				)
			}
		}
	}
	return ctx.Err()
}
