package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"time"
)

type Config struct {
	Server ServerConfig
	Redis  RedisConfig
	Poll   PollConfig
}

type ServerConfig struct {
	Host            string        `env:"SERVER_HOST" env-default:"0.0.0.0"`
	Port            int           `env:"SERVER_PORT" env-default:"8080"`
	ReadTimeout     time.Duration `env:"SERVER_READ_TIMEOUT" env-default:"10s"`
	WriteTimeout    time.Duration `env:"SERVER_WRITE_TIMEOUT" env-default:"10s"`
	ShutdownTimeout time.Duration `env:"SERVER_SHUTDOWN_TIMEOUT" env-default:"5s"`
	BaseURL         string        `env:"BASE_URL" env-default:"http://localhost:8080"`
}

type RedisConfig struct {
	Host     string `env:"REDIS_HOST" env-default:"localhost"`
	Port     int    `env:"REDIS_PORT" env-default:"6379"`
	Password string `env:"REDIS_PASSWORD" env-default:""`
	DB       int    `env:"REDIS_DB" env-default:"0"`
}

type PollConfig struct {
	DefaultTTL time.Duration `env:"POLL_DEFAULT_TTL" env-default:"168h"` // 7 дней
	MaxTTL     time.Duration `env:"POLL_MAX_TTL" env-default:"720h"`     // 30 дней
}

func Load() (*Config, error) {
	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	return &cfg, nil
}

// Address returns Server address as host:port
func (s ServerConfig) Address() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

// Address returns Redis address as host:port
func (r RedisConfig) Address() string {
	return fmt.Sprintf("%s:%d", r.Host, r.Port)
}
