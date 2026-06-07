package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/NicolasPaterno/warden-engine/internal/config"
	"github.com/NicolasPaterno/warden-engine/internal/nats"
	sensorv1 "github.com/NicolasPaterno/warden-proto/gen/go/warden/sensor/v1"
	"github.com/joho/godotenv"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	if err := godotenv.Load(); err != nil {
		slog.Info("no .env file found")
	}
	cfg := config.Load()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	sub, err := nats.NewSubscriber(cfg.NATSUrl)
	if err != nil {
		slog.Error("failed to connect to NATS", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := sub.Close(); err != nil {
			slog.Error("failed to close NATS", "error", err)
		}
	}()

	if err := sub.Subscribe(ctx, func(ctx context.Context, reading *sensorv1.SensorReading) {
		slog.Info("reading received",
			"room", reading.Room,
			"type", reading.Type,
			"value", reading.Value)
	}); err != nil {
		slog.Error("failed to subscribe", "error", err)
		os.Exit(1)
	}

	<-ctx.Done()
	slog.Info("shutting down")
}
