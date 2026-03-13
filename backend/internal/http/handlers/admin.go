package handlers

import (
	"fmt"
	"net/http"
	"net/mail"
	"strconv"
	"strings"
	"time"

	"github.com/darkvoice1/multi-tenant-saas-platform/backend/internal/auth"
	"github.com/darkvoice1/multi-tenant-saas-platform/backend/internal/middleware"
	"github.com/darkvoice1/multi-tenant-saas-platform/backend/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	maxUserNameLen = 80
	minPasswordLen = 8
)

type AdminHandler struct {
	DB *gorm.DB
}

func NewAdminHandler(db *gorm.DB) *AdminHandler {
	return &AdminHandler{DB: db}
}

func (h *AdminHandler) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "admin_ok",
	})
}

type createUserRequest struct {
	Email    string `json:"email" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Role     string `json:"role"`
	Password string `json:"password" binding:"required"`
}

func (h *AdminHandler) CreateUser(c *gin.Context) {
	tenantID, ok := tenantUUIDFromContext(c)
	if !ok {
		return
	}

	var req createUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	req.Name = strings.TrimSpace(req.Name)
	req.Role = strings.ToLower(strings.TrimSpace(req.Role))

	if req.Role == "" {
		req.Role = auth.RoleMember
	}
	if !auth.IsValidRole(req.Role) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role"})
		return
	}
	if req.Name == "" || len(req.Name) > maxUserNameLen {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid name"})
		return
	}
	if _, err := mail.ParseAddress(req.Email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email"})
		return
	}
	if err := validatePassword(req.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var tenant models.Tenant
	if err := h.DB.Select("id", "max_members").Where("id = ?", tenantID).First(&tenant).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "tenant not found"})
		return
	}

	if tenant.MaxMembers > 0 {
		var count int64
		if err := h.DB.Model(&models.User{}).Where("tenant_id = ?", tenantID).Count(&count).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "count failed"})
			return
		}
		if count >= int64(tenant.MaxMembers) {
			c.JSON(http.StatusConflict, gin.H{"error": "member quota exceeded", "code": "MEMBER_QUOTA_EXCEEDED"})
			return
		}
	}

	var exists int64
	if err := h.DB.Model(&models.User{}).Where("tenant_id = ? AND email = ?", tenantID, req.Email).Count(&exists).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "query failed"})
		return
	}
	if exists > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "email already exists"})
		return
	}

	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "hash password failed"})
		return
	}

	user := models.User{
		TenantID:     tenantID,
		Email:        req.Email,
		Name:         req.Name,
		Role:         req.Role,
		Status:       "active",
		PasswordHash: hash,
	}

	if err := h.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "create user failed"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":    user.ID.String(),
		"email": user.Email,
		"name":  user.Name,
		"role":  user.Role,
	})
}

func (h *AdminHandler) ListAuditLogs(c *gin.Context) {
	tenantID, ok := tenantUUIDFromContext(c)
	if !ok {
		return
	}

	limit := 100
	if v := c.Query("limit"); v != "" {
		if parsed, err := parsePositiveInt(v); err == nil {
			if parsed > 0 && parsed <= 200 {
				limit = parsed
			}
		}
	}

	var before time.Time
	if v := c.Query("before"); v != "" {
		if parsed, err := time.Parse(time.RFC3339, v); err == nil {
			before = parsed
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid before"})
			return
		}
	}

	query := h.DB.Where("tenant_id = ?", tenantID)
	if !before.IsZero() {
		query = query.Where("created_at < ?", before)
	}

	var logs []models.AuditLog
	if err := query.Order("created_at desc").Limit(limit).Find(&logs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "query failed"})
		return
	}

	c.JSON(http.StatusOK, logs)
}

func tenantUUIDFromContext(c *gin.Context) (uuid.UUID, bool) {
	tenantIDStr, ok := middleware.TenantID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing tenant id"})
		return uuid.UUID{}, false
	}
	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tenant id"})
		return uuid.UUID{}, false
	}
	return tenantID, true
}

func validatePassword(password string) error {
	if len(password) < minPasswordLen {
		return fmt.Errorf("password too short")
	}
	var hasUpper, hasLower, hasDigit bool
	for _, r := range password {
		switch {
		case r >= 'A' && r <= 'Z':
			hasUpper = true
		case r >= 'a' && r <= 'z':
			hasLower = true
		case r >= '0' && r <= '9':
			hasDigit = true
		}
	}
	if !hasUpper || !hasLower || !hasDigit {
		return fmt.Errorf("password must include upper, lower, digit")
	}
	return nil
}

func parsePositiveInt(input string) (int, error) {
	v, err := strconv.Atoi(input)
	if err != nil || v < 0 {
		return 0, fmt.Errorf("invalid")
	}
	return v, nil
}
