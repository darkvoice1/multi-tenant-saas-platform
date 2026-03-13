package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/darkvoice1/multi-tenant-saas-platform/backend/internal/observability"
	"github.com/gin-gonic/gin"
)

func TestRequestLogger(t *testing.T) {
	gin.SetMode(gin.TestMode)
	observability.InitMetrics()

	r := gin.New()
	r.Use(TraceMiddleware())
	r.Use(RequestLogger())
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if w.Header().Get(TraceIDHeader) == "" {
		t.Fatalf("expected trace header")
	}
}
