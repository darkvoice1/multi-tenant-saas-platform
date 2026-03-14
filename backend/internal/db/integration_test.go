//go:build integration

package db

import (
	"os"
	"testing"
	"time"

	"github.com/darkvoice1/multi-tenant-saas-platform/backend/internal/config"
)

func TestConnectAndMigrateIntegration(t *testing.T) {
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		t.Skip("DATABASE_URL not set")
	}

	cfg := config.Config{
		DatabaseURL:      url,
		DBLogLevel:       "silent",
		SlowSQLThreshold: 200 * time.Millisecond,
	}

	conn, err := Connect(cfg)
	if err != nil {
		t.Fatalf("Connect error: %v", err)
	}

	if err := ApplyMigrations(conn); err != nil {
		t.Fatalf("ApplyMigrations error: %v", err)
	}
}
