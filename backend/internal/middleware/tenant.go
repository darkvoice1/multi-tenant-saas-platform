package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const TenantIDHeader = "X-Tenant-ID"
const tenantIDKey = "tenant_id"

// TenantMiddleware injects tenant_id into context. If required is true,
// missing tenant_id will return 400.
func TenantMiddleware(required bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID := c.GetHeader(TenantIDHeader)
		if tenantID == "" && required {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "missing tenant id",
			})
			c.Abort()
			return
		}

		if tenantID != "" {
			c.Set(tenantIDKey, tenantID)
		}

		c.Next()
	}
}

// TenantID returns tenant_id from gin context.
func TenantID(c *gin.Context) (string, bool) {
	v, ok := c.Get(tenantIDKey)
	if !ok {
		return "", false
	}
	id, ok := v.(string)
	return id, ok
}
