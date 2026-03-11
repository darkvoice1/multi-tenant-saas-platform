package config

import (
	"errors"
	"os"
	"time"
)

// Config holds runtime configuration values.
type Config struct {
	ServerAddr      string
	DatabaseURL     string
	Environment     string
	JWTSecret       string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

// Load reads configuration from environment variables.
func Load() (Config, error) {
	cfg := Config{
		ServerAddr:  getEnv("SERVER_ADDR", ":8080"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
		Environment: getEnv("ENV", "dev"),
		JWTSecret:   os.Getenv("JWT_SECRET"),
	}

	if cfg.DatabaseURL == "" {
		return cfg, errors.New("DATABASE_URL is required")
	}
	if cfg.JWTSecret == "" {
		return cfg, errors.New("JWT_SECRET is required")
	}

	accessTTL, err := time.ParseDuration(getEnv("ACCESS_TOKEN_TTL", "15m"))
	if err != nil {
		return cfg, errors.New("invalid ACCESS_TOKEN_TTL")
	}
	refreshTTL, err := time.ParseDuration(getEnv("REFRESH_TOKEN_TTL", "168h"))
	if err != nil {
		return cfg, errors.New("invalid REFRESH_TOKEN_TTL")
	}
	cfg.AccessTokenTTL = accessTTL
	cfg.RefreshTokenTTL = refreshTTL

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
