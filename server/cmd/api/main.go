package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/sitcon-tw/camp2026-game/internal/app"
)

// @title Camp 2026 Game API
// @version 0.1.0
// @description Backend API for the SITCON Camp 2026 game.
// @BasePath /api
// @schemes http https
// @securityDefinitions.apikey AuthCookieAuth
// @in header
// @name Cookie
// @description User auth cookie. Send as `Cookie: camp2026_auth=<base64url-auth-token>`.
func main() {
	os.Exit(run())
}

func run() int {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	application, err := app.New(ctx)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "startup failed: %v\n", err)
		return 1
	}

	if err := application.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
		application.Log.Error("application stopped with error", "error", err)
		return 1
	}

	return 0
}
