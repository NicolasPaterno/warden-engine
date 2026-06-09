package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"net/http"

	"github.com/NicolasPaterno/warden-engine/internal/config"
	httptransport "github.com/NicolasPaterno/warden-engine/internal/http"
	"github.com/NicolasPaterno/warden-engine/internal/nats"
	"github.com/NicolasPaterno/warden-engine/internal/postgres"
	"github.com/NicolasPaterno/warden-engine/internal/service"
	"github.com/NicolasPaterno/warden-engine/internal/tracing"
	sensorv1 "github.com/NicolasPaterno/warden-proto/gen/go/warden/sensor/v1"
	"github.com/jackc/pgx/v5/pgxpool"
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

	shutdown, err := tracing.Init(ctx, cfg.JaegerEndpoint)
	if err != nil {
		slog.Error("failed to init tracing", "error", err)
		os.Exit(1)
	}
	defer func() {
		flushCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := shutdown(flushCtx); err != nil {
			slog.Error("failed to shutdown tracing", "error", err)
		}
	}()

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

	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	ruleRepo := postgres.NewRuleRepo(pool)
	alertRepo := postgres.NewAlertRepo(pool)

	ruleService := service.NewRuleService(ruleRepo)
	alertService := service.NewAlertService(alertRepo)
	engineService := service.NewEngineService(ruleRepo, alertRepo)

	healthHandler := httptransport.NewHealthHandler(pool, sub)
	router := httptransport.NewRouter(ruleService, alertService, healthHandler)
	server := &http.Server{Addr: cfg.HTTPPort, Handler: router}
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("failed to listen and serve", "error", err)
		}
	}()

	if err := sub.Subscribe(ctx, func(msgCtx context.Context, reading *sensorv1.SensorReading) {
		if err := engineService.Evaluate(msgCtx, reading); err != nil {
			slog.Error("failed to evaluate engine", "error", err)
		}
	}); err != nil {
		slog.Error("failed to subscribe", "error", err)
	}

	<-ctx.Done()
	slog.Info("shutting down")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("failed to shutdown gracefully", "error", err)
	}

}
