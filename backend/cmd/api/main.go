package main

import (
	"context"
	"github.com/AlexeyLars/surway-service/internal/config"
	"github.com/AlexeyLars/surway-service/internal/handler"
	"github.com/AlexeyLars/surway-service/internal/service"
	"github.com/AlexeyLars/surway-service/internal/storage"
	"github.com/redis/go-redis/v9"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/AlexeyLars/surway-service/docs"
)

// @title           Poll Service API
// @version         1.0
// @description     API service for creation and voting in pools
// @termsOfService  http://swagger.io/terms/

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api/v1

// @schemes http https
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

	// Initialize storage, logic and handler layers
	stor := storage.NewRedisStorage(redisClient)
	pollService := service.NewPollService(stor, cfg, logger)
	pollHandler := handler.NewPollHandler(pollService, logger)

	// Setup router and http server
	router := handler.SetupRouter(pollHandler, logger, cfg.Env == "prod")
	server := &http.Server{
		Addr:         cfg.Server.Address(),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	go func() {
		logger.Info("Starting http server", slog.String("address", server.Addr))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("failed to start server", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	ctx, cancel = context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("server forced to shutdown", slog.String("error", err.Error()))
	}

	// Close Redis connection
	if err := stor.Close(); err != nil {
		logger.Error("failed to close redis connection", slog.String("error", err.Error()))
	}

	logger.Info("Server stopped")
}
