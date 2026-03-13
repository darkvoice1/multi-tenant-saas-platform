package handlers

import (
	"net/http"

	"github.com/darkvoice1/multi-tenant-saas-platform/backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

func TenantEcho(c *gin.Context) {
	tenantID, _ := middleware.TenantID(c)
	c.JSON(http.StatusOK, gin.H{
		"tenant_id": tenantID,
	})
}
