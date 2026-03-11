package middleware

import (
	"net/http"

	"github.com/darkvoice1/multi-tenant-saas-platform/backend/internal/auth"
	"github.com/gin-gonic/gin"
)

func RequireRoles(roles ...string) gin.HandlerFunc {
	allowed := map[string]struct{}{}
	for _, r := range roles {
		allowed[r] = struct{}{}
	}
	return func(c *gin.Context) {
		role, ok := c.Get(contextRoleKey)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing role"})
			c.Abort()
			return
		}
		if _, ok := allowed[role.(string)]; !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func RequirePermission(perm auth.Permission) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, ok := c.Get(contextRoleKey)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing role"})
			c.Abort()
			return
		}
		if !auth.IsAllowed(role.(string), perm) {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			c.Abort()
			return
		}
		c.Next()
	}
}
