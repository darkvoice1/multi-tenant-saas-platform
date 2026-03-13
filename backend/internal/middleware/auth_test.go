package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/darkvoice1/multi-tenant-saas-platform/backend/internal/auth"
	"github.com/gin-gonic/gin"
)

func TestAuthMiddlewareMissingHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(AuthMiddleware("secret"))
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestAuthMiddlewareInvalidHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(AuthMiddleware("secret"))
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Token abc")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestAuthMiddlewareTenantMismatch(t *testing.T) {
	gin.SetMode(gin.TestMode)
	token, err := auth.CreateAccessToken("secret", time.Minute, "user1", "tenant-a", auth.RoleAdmin)
	if err != nil {
		t.Fatalf("CreateAccessToken error: %v", err)
	}

	r := gin.New()
	r.Use(TenantMiddleware(true))
	r.Use(AuthMiddleware("secret"))
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set(TenantIDHeader, "tenant-b")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", w.Code)
	}
}

func TestAuthMiddlewareSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	token, err := auth.CreateAccessToken("secret", time.Minute, "user1", "tenant-a", auth.RoleMember)
	if err != nil {
		t.Fatalf("CreateAccessToken error: %v", err)
	}

	r := gin.New()
	r.Use(TenantMiddleware(true))
	r.Use(AuthMiddleware("secret"))
	r.GET("/", func(c *gin.Context) {
		uid, role, tid, ok := CurrentUser(c)
		if !ok || uid == "" || role == "" || tid == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing user"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set(TenantIDHeader, "tenant-a")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}
