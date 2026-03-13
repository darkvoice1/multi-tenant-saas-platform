package config

import (
	"testing"
	"time"
)

func TestLoadMissingRequired(t *testing.T) {
	t.Setenv("DATABASE_URL", "")
	t.Setenv("JWT_SECRET", "")
	if _, err := Load(); err == nil {
		t.Fatalf("expected error for missing DATABASE_URL and JWT_SECRET")
	}
}

func TestLoadWeakSecretInProd(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://user:pass@localhost/db")
	t.Setenv("JWT_SECRET", "dev_change_me")
	t.Setenv("ENV", "prod")
	if _, err := Load(); err == nil {
		t.Fatalf("expected weak secret error in prod")
	}
}

func TestLoadInvalidDurations(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://user:pass@localhost/db")
	t.Setenv("JWT_SECRET", "strong_secret_123")
	t.Setenv("ACCESS_TOKEN_TTL", "notaduration")
	if _, err := Load(); err == nil {
		t.Fatalf("expected invalid ACCESS_TOKEN_TTL error")
	}
}

func TestLoadDefaultsAndBools(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://user:pass@localhost/db")
	t.Setenv("JWT_SECRET", "strong_secret_123")
	t.Setenv("S3_USE_SSL", "true")
	t.Setenv("S3_PATH_STYLE", "false")
	t.Setenv("ACCESS_TOKEN_TTL", "10m")
	t.Setenv("REFRESH_TOKEN_TTL", "1h")
	t.Setenv("SLOW_SQL_THRESHOLD", "150ms")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load error: %v", err)
	}
	if cfg.ServerAddr == "" || cfg.StorageDir == "" {
		t.Fatalf("expected defaults to be set")
	}
	if !cfg.S3UseSSL {
		t.Fatalf("expected S3UseSSL true")
	}
	if cfg.S3PathStyle {
		t.Fatalf("expected S3PathStyle false")
	}
	if cfg.AccessTokenTTL != 10*time.Minute || cfg.RefreshTokenTTL != time.Hour {
		t.Fatalf("unexpected ttl values")
	}
	if cfg.SlowSQLThreshold != 150*time.Millisecond {
		t.Fatalf("unexpected slow sql threshold")
	}
}
