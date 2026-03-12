package db

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

// ApplyMigrations runs SQL migrations in order and records applied versions.
func ApplyMigrations(database *gorm.DB) error {
	dir, err := resolveMigrationDir()
	if err != nil {
		return err
	}

	if err := database.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version bigint PRIMARY KEY,
			dirty boolean NOT NULL DEFAULT false
		)
	`).Error; err != nil {
		return fmt.Errorf("create migrations table: %w", err)
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("read migrations dir: %w", err)
	}

	var files []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasSuffix(name, ".up.sql") {
			files = append(files, name)
		}
	}
	sort.Strings(files)

	for _, name := range files {
		version, ok := parseNumericVersion(name)
		if !ok {
			continue
		}
		var count int64
		if err := database.Raw("SELECT count(1) FROM schema_migrations WHERE version = ?", version).
			Scan(&count).Error; err != nil {
			return fmt.Errorf("check migration %s: %w", name, err)
		}
		if count > 0 {
			continue
		}

		path := filepath.Join(dir, name)
		bytes, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read migration %s: %w", name, err)
		}
		sql := strings.TrimPrefix(string(bytes), "\ufeff")
		if strings.TrimSpace(sql) == "" {
			continue
		}
		if err := database.Exec(sql).Error; err != nil {
			return fmt.Errorf("apply migration %s: %w", name, err)
		}
		if err := database.Exec("INSERT INTO schema_migrations (version, dirty) VALUES (?, false)", version).Error; err != nil {
			return fmt.Errorf("record migration %s: %w", name, err)
		}
	}

	return nil
}

func resolveMigrationDir() (string, error) {
	candidates := []string{
		"migrations",
		filepath.Join("backend", "migrations"),
	}
	for _, dir := range candidates {
		if info, err := os.Stat(dir); err == nil && info.IsDir() {
			return dir, nil
		}
	}
	return "", fmt.Errorf("migrations directory not found")
}

func parseNumericVersion(filename string) (int64, bool) {
	base := strings.TrimSuffix(filename, ".up.sql")
	prefix, _, ok := strings.Cut(base, "_")
	if !ok {
		return 0, false
	}
	n, err := strconv.ParseInt(prefix, 10, 64)
	if err != nil {
		return 0, false
	}
	return n, true
}
