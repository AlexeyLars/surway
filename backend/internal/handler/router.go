package handler

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

// SetupRouter setups router with middleware and routes
func SetupRouter(handler *PollHandler, logger *slog.Logger, releaseMode bool) *gin.Engine {
	if releaseMode {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Middleware
	router.Use(gin.Recovery())
	router.Use(LoggerMiddleware(logger))
	router.Use(CORSMiddleware())

	// Health check
	router.GET("/health", handler.HealthCheck)

	// API routes
	v1 := router.Group("/api/v1")
	{
		polls := v1.Group("/polls")
		{
			polls.POST("", handler.CreatePoll)
			polls.POST("/:id/vote", handler.Vote)
			polls.GET("/:id/results", handler.GetResults)
		}
	}

	return router
}

// LoggerMiddleware logs HTTP requests
func LoggerMiddleware(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		duration := time.Since(start)

		logger.InfoContext(c.Request.Context(), "http request",
			slog.String("method", c.Request.Method),
			slog.String("path", path),
			slog.String("query", query),
			slog.Int("status", c.Writer.Status()),
			slog.Duration("duration", duration),
			slog.String("client_ip", c.ClientIP()),
			slog.String("user_agent", c.Request.UserAgent()),
		)
	}
}

// CORSMiddleware add CORS headers
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
