package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	TraceIDHeader = "X-Trace-ID"
	traceIDKey    = "trace_id"
)

func TraceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := strings.TrimSpace(c.GetHeader(TraceIDHeader))
		if traceID == "" {
			traceID = uuid.New().String()
		}
		c.Set(traceIDKey, traceID)
		c.Writer.Header().Set(TraceIDHeader, traceID)
		c.Next()
	}
}

func TraceID(c *gin.Context) (string, bool) {
	v, ok := c.Get(traceIDKey)
	if !ok {
		return "", false
	}
	id, ok := v.(string)
	return id, ok
}
