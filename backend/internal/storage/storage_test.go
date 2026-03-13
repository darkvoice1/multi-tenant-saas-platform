package storage

import (
	"testing"

	"github.com/darkvoice1/multi-tenant-saas-platform/backend/internal/config"
)

func TestNewStorageLocal(t *testing.T) {
	cfg := config.Config{StorageBackend: "local", StorageDir: t.TempDir()}
	store, err := NewStorage(cfg)
	if err != nil {
		t.Fatalf("NewStorage error: %v", err)
	}
	if store == nil {
		t.Fatalf("expected storage")
	}
}

func TestNewStorageUnknown(t *testing.T) {
	cfg := config.Config{StorageBackend: "unknown"}
	if _, err := NewStorage(cfg); err == nil {
		t.Fatalf("expected error")
	}
}
