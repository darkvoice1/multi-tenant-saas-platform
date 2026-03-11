package config

import (
	"errors"
	"os"
)

// Config holds runtime configuration values.
type Config struct {
	ServerAddr string
	DatabaseURL string
	Environment string
}

// Load reads configuration from environment variables.
func Load() (Config, error) {
	cfg := Config{
		ServerAddr:  getEnv("SERVER_ADDR", ":8080"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
		Environment: getEnv("ENV", "dev"),
	}

	if cfg.DatabaseURL == "" {
		return cfg, errors.New("DATABASE_URL is required")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
