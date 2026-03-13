package observability

import (
	"context"
	"testing"

	"github.com/darkvoice1/multi-tenant-saas-platform/backend/internal/config"
)

func TestInitTracerNone(t *testing.T) {
	cfg := config.Config{OTelExporter: "none", OTelServiceName: "svc"}
	shutdown, err := InitTracer(cfg)
	if err != nil {
		t.Fatalf("InitTracer error: %v", err)
	}
	if shutdown == nil {
		t.Fatalf("expected shutdown func")
	}
	_ = shutdown(context.Background())
}

func TestInitTracerStdout(t *testing.T) {
	cfg := config.Config{OTelExporter: "stdout", OTelServiceName: "svc"}
	shutdown, err := InitTracer(cfg)
	if err != nil {
		t.Fatalf("InitTracer error: %v", err)
	}
	if shutdown == nil {
		t.Fatalf("expected shutdown func")
	}
	_ = shutdown(context.Background())
}

func TestInitTracerMissingEndpoint(t *testing.T) {
	cfg := config.Config{OTelExporter: "otlp", OTelServiceName: "svc"}
	if _, err := InitTracer(cfg); err == nil {
		t.Fatalf("expected error for missing endpoint")
	}
}
