package db

import (
	"fmt"

	"github.com/darkvoice1/multi-tenant-saas-platform/backend/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Connect opens a database connection using GORM.
func Connect(cfg config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get sql db: %w", err)
	}
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("ping db: %w", err)
	}

	return db, nil
}
