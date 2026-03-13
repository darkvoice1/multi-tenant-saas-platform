package config

import (
	"errors"
	"os"
	"strings"
	"time"
)

// Config holds runtime configuration values.
type Config struct {
	ServerAddr       string
	DatabaseURL      string
	Environment      string
	DBLogLevel       string
	SlowSQLThreshold time.Duration
	JWTSecret        string
	AccessTokenTTL   time.Duration
	RefreshTokenTTL  time.Duration
	StorageDir       string
	StorageBackend   string
	S3Endpoint       string
	S3Region         string
	S3Bucket         string
	S3AccessKey      string
	S3SecretKey      string
	S3UseSSL         bool
	S3PathStyle      bool
}

// Load reads configuration from environment variables.
func Load() (Config, error) {
	cfg := Config{
		ServerAddr:     getEnv("SERVER_ADDR", ":8080"),
		DatabaseURL:    os.Getenv("DATABASE_URL"),
		Environment:    getEnv("ENV", "dev"),
		DBLogLevel:     getEnv("DB_LOG_LEVEL", "warn"),
		JWTSecret:      os.Getenv("JWT_SECRET"),
		StorageDir:     getEnv("STORAGE_DIR", "storage"),
		StorageBackend: getEnv("STORAGE_BACKEND", "local"),
		S3Endpoint:     os.Getenv("S3_ENDPOINT"),
		S3Region:       getEnv("S3_REGION", "us-east-1"),
		S3Bucket:       getEnv("S3_BUCKET", "saas-platform"),
		S3AccessKey:    os.Getenv("S3_ACCESS_KEY"),
		S3SecretKey:    os.Getenv("S3_SECRET_KEY"),
		S3UseSSL:       getEnvBool("S3_USE_SSL", false),
		S3PathStyle:    getEnvBool("S3_PATH_STYLE", true),
	}

	if cfg.DatabaseURL == "" {
		return cfg, errors.New("DATABASE_URL is required")
	}
	if cfg.JWTSecret == "" {
		return cfg, errors.New("JWT_SECRET is required")
	}
	if cfg.Environment != "dev" {
		if len(cfg.JWTSecret) < 16 || strings.EqualFold(cfg.JWTSecret, "dev_change_me") {
			return cfg, errors.New("JWT_SECRET is too weak")
		}
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

	slowThreshold, err := time.ParseDuration(getEnv("SLOW_SQL_THRESHOLD", "200ms"))
	if err != nil {
		return cfg, errors.New("invalid SLOW_SQL_THRESHOLD")
	}
	cfg.SlowSQLThreshold = slowThreshold

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	if v := os.Getenv(key); v != "" {
		switch v {
		case "1", "true", "TRUE", "True", "yes", "YES", "Yes", "on", "ON", "On":
			return true
		case "0", "false", "FALSE", "False", "no", "NO", "No", "off", "OFF", "Off":
			return false
		}
	}
	return fallback
}
