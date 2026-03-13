package middleware

import (
	"encoding/json"
	"log"
	"time"

	"github.com/darkvoice1/multi-tenant-saas-platform/backend/internal/observability"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
)

type requestLogEntry struct {
	TraceID   string `json:"trace_id"`
	TenantID  string `json:"tenant_id,omitempty"`
	UserID    string `json:"user_id,omitempty"`
	Method    string `json:"method"`
	Path      string `json:"path"`
	Status    int    `json:"status"`
	LatencyMs int64  `json:"latency_ms"`
	ClientIP  string `json:"client_ip"`
	UserAgent string `json:"user_agent"`
	CreatedAt string `json:"created_at"`
}

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		status := c.Writer.Status()
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}

		traceID := ""
		if span := trace.SpanFromContext(c.Request.Context()); span.SpanContext().IsValid() {
			traceID = span.SpanContext().TraceID().String()
		}
		if traceID == "" {
			if v, ok := TraceID(c); ok {
				traceID = v
			}
		}

		tenantID, _ := TenantID(c)
		userID, _, _, _ := CurrentUser(c)
		latency := time.Since(start)

		observability.RecordRequest(c.Request.Method, path, status, latency)

		entry := requestLogEntry{
			TraceID:   traceID,
			TenantID:  tenantID,
			UserID:    userID,
			Method:    c.Request.Method,
			Path:      path,
			Status:    status,
			LatencyMs: latency.Milliseconds(),
			ClientIP:  c.ClientIP(),
			UserAgent: c.Request.UserAgent(),
			CreatedAt: time.Now().Format(time.RFC3339),
		}
		if data, err := json.Marshal(entry); err == nil {
			log.Print(string(data))
		}
	}
}
