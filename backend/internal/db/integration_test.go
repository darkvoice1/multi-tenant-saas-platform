//go:build integration

package db

import (
	"os"
	"path/filepath"
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

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd error: %v", err)
	}
	backendDir := filepath.Clean(filepath.Join(wd, "..", ".."))
	if err := os.Chdir(backendDir); err != nil {
		t.Fatalf("chdir backend error: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(wd)
	})

	conn, err := Connect(cfg)
	if err != nil {
		t.Fatalf("Connect error: %v", err)
	}

	if err := ApplyMigrations(conn); err != nil {
		t.Fatalf("ApplyMigrations error: %v", err)
	}
}