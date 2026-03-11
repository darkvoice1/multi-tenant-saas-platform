package middleware

import (
	"net/http"
	"strings"

	"github.com/darkvoice1/multi-tenant-saas-platform/backend/internal/auth"
	"github.com/gin-gonic/gin"
)

const (
	contextUserIDKey   = "user_id"
	contextRoleKey     = "user_role"
	contextTenantIDKey = "token_tenant_id"
)

func AuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			c.Abort()
			return
		}
		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header"})
			c.Abort()
			return
		}

		claims, err := auth.ParseToken(secret, parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		c.Set(contextUserIDKey, claims.UserID)
		c.Set(contextRoleKey, claims.Role)
		c.Set(contextTenantIDKey, claims.TenantID)

		if tenantID, ok := TenantID(c); ok && tenantID != claims.TenantID {
			c.JSON(http.StatusForbidden, gin.H{"error": "tenant mismatch"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func CurrentUser(c *gin.Context) (string, string, string, bool) {
	userID, ok := c.Get(contextUserIDKey)
	if !ok {
		return "", "", "", false
	}
	role, ok := c.Get(contextRoleKey)
	if !ok {
		return "", "", "", false
	}
	tenantID, _ := c.Get(contextTenantIDKey)
	return userID.(string), role.(string), tenantID.(string), true
}
