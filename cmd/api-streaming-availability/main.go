// Command api-streaming-availability serves segment availability for a CDN origin.
package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/LoyalMinion/api-streaming-availability/internal/config"
	"github.com/LoyalMinion/api-streaming-availability/internal/httpapi"
	"github.com/LoyalMinion/api-streaming-availability/internal/server"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	if err := run(logger); err != nil {
		logger.Error("run", "err", err)
		os.Exit(1)
	}
}

func run(logger *slog.Logger) error {
	cfg, err := config.Load(os.Getenv)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	api := httpapi.New(logger, time.Now(), time.Now)
	if err := server.New(cfg, api.Handler(), logger).ListenAndServe(ctx); err != nil {
		return fmt.Errorf("run server: %w", err)
	}

	return nil
}
