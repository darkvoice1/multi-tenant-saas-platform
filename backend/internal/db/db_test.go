package db

import (
	"testing"

	"gorm.io/gorm/logger"
)

func TestParseLogLevel(t *testing.T) {
	if parseLogLevel("silent") != logger.Silent {
		t.Fatalf("expected silent")
	}
	if parseLogLevel("error") != logger.Error {
		t.Fatalf("expected error")
	}
	if parseLogLevel("info") != logger.Info {
		t.Fatalf("expected info")
	}
	if parseLogLevel("warn") != logger.Warn {
		t.Fatalf("expected warn")
	}
	if parseLogLevel("unknown") != logger.Warn {
		t.Fatalf("expected default warn")
	}
}
