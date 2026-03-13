package observability

import (
	"context"
	"fmt"
	"strings"

	"github.com/darkvoice1/multi-tenant-saas-platform/backend/internal/config"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv/v1.26.0"
)

func InitTracer(cfg config.Config) (func(context.Context) error, error) {
	exporter := strings.ToLower(strings.TrimSpace(cfg.OTelExporter))
	if exporter == "" || exporter == "none" {
		return func(context.Context) error { return nil }, nil
	}

	res, err := resource.New(context.Background(),
		resource.WithAttributes(
			semconv.ServiceName(cfg.OTelServiceName),
		),
	)
	if err != nil {
		return nil, err
	}

	var spanExporter trace.SpanExporter
	if exporter == "stdout" {
		spanExporter, err = stdouttrace.New(stdouttrace.WithPrettyPrint())
		if err != nil {
			return nil, err
		}
	} else if exporter == "otlp" {
		endpoint := strings.TrimSpace(cfg.OTelEndpoint)
		if endpoint == "" {
			return nil, fmt.Errorf("OTEL_EXPORTER_OTLP_ENDPOINT is required")
		}
		spanExporter, err = otlptracegrpc.New(
			context.Background(),
			otlptracegrpc.WithEndpoint(endpoint),
			otlptracegrpc.WithInsecure(),
		)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("unsupported OTEL_EXPORTER: %s", exporter)
	}

	provider := trace.NewTracerProvider(
		trace.WithBatcher(spanExporter),
		trace.WithResource(res),
	)
	otel.SetTracerProvider(provider)

	return provider.Shutdown, nil
}

func OTelGinMiddleware(cfg config.Config) gin.HandlerFunc {
	service := cfg.OTelServiceName
	if strings.TrimSpace(service) == "" {
		service = "saas-platform-api"
	}
	return otelgin.Middleware(service)
}
