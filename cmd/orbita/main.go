package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/vxssroott/ORBITA/internal/events"
	"github.com/vxssroott/ORBITA/internal/rules"
	"github.com/vxssroott/ORBITA/internal/state"
	"github.com/vxssroott/ORBITA/internal/telemetry"
	"github.com/vxssroott/ORBITA/services/api"
)

const version = "0.2.0-dev"

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	logger.Info(
		"starting ORBITA",
		"version", version,
		"service", "orbita-core",
	)

	telemetryStore := telemetry.NewStore()
	stateEngine := state.NewEngine()
	eventEngine := events.NewEngine()

	_ = eventEngine

	ruleEngine := rules.NewEngine(nil)
	_ = ruleEngine

	apiServer := api.NewServer(
		telemetryStore,
		stateEngine,
	)

	httpServer := &http.Server{
		Addr:              "127.0.0.1:8080",
		Handler:           apiServer.Handler(),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		logger.Info("API listening", "address", httpServer.Addr)

		if err := httpServer.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {
			logger.Error("API server stopped", "error", err)
			os.Exit(1)
		}
	}()

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	<-ctx.Done()

	logger.Info("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		logger.Error("graceful shutdown failed", "error", err)
		os.Exit(1)
	}

	fmt.Println("ORBITA shutdown complete")
}
