package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/sitcon-tw/camp2026-game/internal/config"
	"github.com/sitcon-tw/camp2026-game/internal/migration"
)

func main() {
	os.Exit(run())
}

func run() int {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	args := os.Args[1:]
	command := "up"
	if len(args) > 0 {
		command = args[0]
		args = args[1:]
	}

	dir := migration.DefaultDir()
	if command == "new" {
		if len(args) != 1 {
			usage()
			return 2
		}
		path, err := migration.CreateFile(dir, args[0])
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "create migration failed: %v\n", err)
			return 1
		}
		fmt.Printf("created %s\n", path)
		return 0
	}

	cfg, err := config.Load()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "load config failed: %v\n", err)
		return 1
	}

	conn, err := migration.Open(ctx, cfg.DatabaseURL)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "migration startup failed: %v\n", err)
		return 1
	}
	defer conn.Close(context.Background())

	runner := migration.Runner{
		Conn: conn,
		Dir:  dir,
		Out:  os.Stdout,
	}

	switch command {
	case "up":
		err = runner.Up(ctx)
	case "down":
		err = runner.Down(ctx)
	case "status":
		err = runner.Status(ctx)
	default:
		usage()
		return 2
	}

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "migration failed: %v\n", err)
		return 1
	}
	return 0
}

func usage() {
	_, _ = fmt.Fprintln(os.Stderr, `usage:
  go run ./cmd/migrate [up]
  go run ./cmd/migrate down
  go run ./cmd/migrate status
  go run ./cmd/migrate new <name>`)
}
