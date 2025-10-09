package main

import (
	"github.com/AlexeyLars/surway-service/internal/config"
	"log/slog"
	"os"
)

func main() {
	// Configure logging
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	// Load cfg
	cfg, err := config.Load()
	if err != nil {
		logger.Error("failed to load config", slog.String("error", err.Error()))
		os.Exit(1)
	}

	logger.Info("Starting service",
		slog.String("server_address", cfg.Server.Address()),
	)

}
