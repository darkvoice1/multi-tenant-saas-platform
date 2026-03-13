package db

import (
	"fmt"
	"os"
	"strings"

	"github.com/darkvoice1/multi-tenant-saas-platform/backend/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Connect opens a database connection using GORM.
func Connect(cfg config.Config) (*gorm.DB, error) {
	gormLogger := logger.New(
		logWriter{},
		logger.Config{
			SlowThreshold:             cfg.SlowSQLThreshold,
			LogLevel:                  parseLogLevel(cfg.DBLogLevel),
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{Logger: gormLogger})
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

type logWriter struct{}

func (logWriter) Printf(format string, args ...any) {
	_, _ = fmt.Fprintf(os.Stdout, format, args...)
}

func parseLogLevel(level string) logger.LogLevel {
	switch strings.ToLower(level) {
	case "silent":
		return logger.Silent
	case "error":
		return logger.Error
	case "info":
		return logger.Info
	default:
		return logger.Warn
	}
}
