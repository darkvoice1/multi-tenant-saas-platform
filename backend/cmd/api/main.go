package main

import (
	"log"
	"os"

	"github.com/darkvoice1/multi-tenant-saas-platform/backend/internal/config"
	"github.com/darkvoice1/multi-tenant-saas-platform/backend/internal/db"
	"github.com/darkvoice1/multi-tenant-saas-platform/backend/internal/http"
	"github.com/darkvoice1/multi-tenant-saas-platform/backend/internal/storage"
)

func main() {
	if _, err := os.Stat(".env"); err == nil {
		_ = config.LoadDotEnv(".env")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	database, err := db.Connect(cfg)
	if err != nil {
		log.Fatalf("db error: %v", err)
	}

	if cfg.Environment == "dev" {
		if err := db.ApplyMigrations(database); err != nil {
			log.Fatalf("migration error: %v", err)
		}
	}

	if cfg.StorageBackend == "" || cfg.StorageBackend == "local" {
		if err := os.MkdirAll(cfg.StorageDir, 0o755); err != nil {
			log.Fatalf("storage dir error: %v", err)
		}
	}

	store, err := storage.NewStorage(cfg)
	if err != nil {
		log.Fatalf("storage error: %v", err)
	}

	router := http.NewRouter(database, cfg, store)
	if err := router.Run(cfg.ServerAddr); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
