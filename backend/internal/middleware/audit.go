package middleware

import (
	"strings"

	"github.com/darkvoice1/multi-tenant-saas-platform/backend/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func AuditLogger(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		tenantIDStr, ok := TenantID(c)
		if !ok {
			return
		}
		userID, _, _, ok := CurrentUser(c)
		if !ok {
			return
		}

		tenantID, err := uuid.Parse(tenantIDStr)
		if err != nil {
			return
		}
		uid, err := uuid.Parse(userID)
		if err != nil {
			return
		}

		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}

		action, resource := deriveAuditAction(c.Request.Method, path)
		var resourceID *uuid.UUID
		if id := c.Param("id"); id != "" {
			if rid, err := uuid.Parse(id); err == nil {
				resourceID = &rid
			}
		}

		entry := models.AuditLog{
			TenantID:   tenantID,
			UserID:     uid,
			Action:     action,
			Resource:   resource,
			ResourceID: resourceID,
			Method:     c.Request.Method,
			Path:       path,
			StatusCode: c.Writer.Status(),
			IP:         c.ClientIP(),
			UserAgent:  c.Request.UserAgent(),
		}

		_ = db.Create(&entry).Error
	}
}

func deriveAuditAction(method, path string) (string, string) {
	resource := "unknown"
	switch {
	case strings.Contains(path, "/projects"):
		resource = "project"
	case strings.Contains(path, "/tasks"):
		resource = "task"
	case strings.Contains(path, "/comments"):
		resource = "comment"
	case strings.Contains(path, "/attachments"):
		resource = "attachment"
	case strings.Contains(path, "/notifications"):
		resource = "notification"
	case strings.Contains(path, "/dashboard"):
		resource = "dashboard"
	case strings.Contains(path, "/admin"):
		resource = "admin"
	case strings.Contains(path, "/tenant"):
		resource = "tenant"
	}

	switch {
	case strings.Contains(path, "/approve"):
		return "approve", resource
	case strings.Contains(path, "/status"):
		return "status_change", resource
	case strings.Contains(path, "/read"):
		return "mark_read", resource
	}

	switch method {
	case "GET":
		return "read", resource
	case "POST":
		return "create", resource
	case "PUT":
		return "update", resource
	case "DELETE":
		return "delete", resource
	default:
		return strings.ToLower(method), resource
	}
}
