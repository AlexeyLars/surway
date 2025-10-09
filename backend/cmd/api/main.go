package main

import (
	"context"
	"github.com/AlexeyLars/surway-service/internal/config"
	"github.com/redis/go-redis/v9"
	"log/slog"
	"os"
	"time"
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

	// Connect Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Address(),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	// Redis connection check
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := redisClient.Ping(ctx).Err(); err != nil {
		logger.Error("failed to connect to redis",
			slog.String("address", cfg.Redis.Address()),
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}

	logger.Info("Connected to Redis", slog.String("redis_address", cfg.Redis.Address()))
}
