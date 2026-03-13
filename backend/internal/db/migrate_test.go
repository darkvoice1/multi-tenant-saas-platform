package db

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseNumericVersion(t *testing.T) {
	if v, ok := parseNumericVersion("202401010101_init.up.sql"); !ok || v != 202401010101 {
		t.Fatalf("unexpected version parse")
	}
	if _, ok := parseNumericVersion("badfile.up.sql"); ok {
		t.Fatalf("expected parse to fail")
	}
}

func TestResolveMigrationDir(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd error: %v", err)
	}
	defer func() {
		_ = os.Chdir(wd)
	}()

	backendDir := filepath.Clean(filepath.Join(wd, "..", ".."))
	if err := os.Chdir(backendDir); err != nil {
		t.Fatalf("Chdir error: %v", err)
	}

	dir, err := resolveMigrationDir()
	if err != nil {
		t.Fatalf("resolveMigrationDir error: %v", err)
	}
	if filepath.Base(dir) != "migrations" {
		t.Fatalf("expected migrations dir, got %s", dir)
	}
}
