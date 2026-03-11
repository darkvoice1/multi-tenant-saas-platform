package main

import (
	"log"

	"github.com/darkvoice1/multi-tenant-saas-platform/backend/internal/config"
	"github.com/darkvoice1/multi-tenant-saas-platform/backend/internal/db"
	"github.com/darkvoice1/multi-tenant-saas-platform/backend/internal/http"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	database, err := db.Connect(cfg)
	if err != nil {
		log.Fatalf("db error: %v", err)
	}

	router := http.NewRouter(database)
	if err := router.Run(cfg.ServerAddr); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
