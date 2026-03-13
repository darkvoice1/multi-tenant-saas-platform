package handlers

import (
	"github.com/darkvoice1/multi-tenant-saas-platform/backend/internal/auth"
	"github.com/darkvoice1/multi-tenant-saas-platform/backend/internal/config"
	"github.com/darkvoice1/multi-tenant-saas-platform/backend/internal/middleware"
	"github.com/darkvoice1/multi-tenant-saas-platform/backend/internal/models"
	"github.com/darkvoice1/multi-tenant-saas-platform/backend/internal/observability"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"io"
	"net/http"
	"time"
)

type AuthHandler struct {
	DB     *gorm.DB
	Config config.Config
}

func NewAuthHandler(db *gorm.DB, cfg config.Config) *AuthHandler {
	return &AuthHandler{DB: db, Config: cfg}
}

type loginRequest struct {
	TenantID string `json:"tenant_id" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type bootstrapRequest struct {
	TenantName    string `json:"tenant_name"`
	TenantSlug    string `json:"tenant_slug"`
	AdminEmail    string `json:"admin_email"`
	AdminPassword string `json:"admin_password"`
	AdminName     string `json:"admin_name"`
}

func issueTokens(db *gorm.DB, cfg config.Config, user models.User) (string, string, int, error) {
	accessToken, err := auth.CreateAccessToken(cfg.JWTSecret, cfg.AccessTokenTTL, user.ID.String(), user.TenantID.String(), user.Role)
	if err != nil {
		return "", "", 0, err
	}
	refreshToken, refreshHash, err := auth.GenerateRefreshToken()
	if err != nil {
		return "", "", 0, err
	}
	refresh := models.RefreshToken{
		TenantID:  user.TenantID,
		UserID:    user.ID,
		TokenHash: refreshHash,
		ExpiresAt: time.Now().Add(cfg.RefreshTokenTTL),
	}
	if err := db.Create(&refresh).Error; err != nil {
		return "", "", 0, err
	}

	return accessToken, refreshToken, int(cfg.AccessTokenTTL.Seconds()), nil
}

func (h *AuthHandler) Bootstrap(c *gin.Context) {
	if h.Config.Environment != "dev" {
		c.JSON(http.StatusForbidden, gin.H{"error": "bootstrap disabled"})
		return
	}

	var req bootstrapRequest
	if err := c.ShouldBindJSON(&req); err != nil && err != io.EOF {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if req.TenantName == "" {
		req.TenantName = "Demo Tenant"
	}
	if req.TenantSlug == "" {
		req.TenantSlug = "demo"
	}
	if req.AdminEmail == "" {
		req.AdminEmail = "admin@example.com"
	}
	if req.AdminPassword == "" {
		req.AdminPassword = "Admin123!"
	}
	if req.AdminName == "" {
		req.AdminName = "Admin"
	}
	if err := validatePassword(req.AdminPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var count int64
	if err := h.DB.Model(&models.Tenant{}).Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "already initialized"})
		return
	}

	tenant := models.Tenant{
		Name:   req.TenantName,
		Slug:   req.TenantSlug,
		Status: "active",
	}
	if err := h.DB.Create(&tenant).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "create tenant failed"})
		return
	}

	hash, err := auth.HashPassword(req.AdminPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "hash password failed"})
		return
	}
	user := models.User{
		TenantID:     tenant.ID,
		Email:        req.AdminEmail,
		Name:         req.AdminName,
		Role:         auth.RoleAdmin,
		Status:       "active",
		PasswordHash: hash,
	}
	if err := h.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "create user failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tenant_id": tenant.ID.String(),
		"tenant":    tenant.Name,
		"email":     user.Email,
		"password":  req.AdminPassword,
		"user_id":   user.ID.String(),
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	tenantID, err := uuid.Parse(req.TenantID)
	if err != nil {
		observability.RecordLogin(false)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tenant_id"})
		return
	}

	var user models.User
	if err := h.DB.Where("tenant_id = ? AND email = ?", tenantID, req.Email).First(&user).Error; err != nil {
		observability.RecordLogin(false)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	if user.Status != "active" {
		observability.RecordLogin(false)
		c.JSON(http.StatusForbidden, gin.H{"error": "user disabled"})
		return
	}
	if user.PasswordHash == "" || !auth.CheckPassword(user.PasswordHash, req.Password) {
		observability.RecordLogin(false)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	accessToken, refreshToken, expiresIn, err := issueTokens(h.DB, h.Config, user)
	if err != nil {
		observability.RecordLogin(false)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token error"})
		return
	}

	now := time.Now()
	h.DB.Model(&user).Updates(map[string]interface{}{
		"last_login_at": &now,
	})

	observability.RecordLogin(true)

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"expires_in":    expiresIn,
		"user": gin.H{
			"id":        user.ID.String(),
			"tenant_id": user.TenantID.String(),
			"email":     user.Email,
			"name":      user.Name,
			"role":      user.Role,
		},
	})
}
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req refreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	hash := auth.HashRefreshToken(req.RefreshToken)
	var token models.RefreshToken
	if err := h.DB.Where("token_hash = ?", hash).First(&token).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}
	if token.RevokedAt != nil || token.ExpiresAt.Before(time.Now()) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "expired refresh token"})
		return
	}

	var user models.User
	if err := h.DB.First(&user, "id = ?", token.UserID).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}
	if user.Status != "active" {
		c.JSON(http.StatusForbidden, gin.H{"error": "user disabled"})
		return
	}

	now := time.Now()
	h.DB.Model(&token).Updates(map[string]interface{}{
		"revoked_at": &now,
	})

	accessToken, refreshToken, expiresIn, err := issueTokens(h.DB, h.Config, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"expires_in":    expiresIn,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	var req refreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	hash := auth.HashRefreshToken(req.RefreshToken)
	var token models.RefreshToken
	if err := h.DB.Where("token_hash = ?", hash).First(&token).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
		return
	}
	if token.RevokedAt == nil {
		now := time.Now()
		h.DB.Model(&token).Updates(map[string]interface{}{
			"revoked_at": &now,
		})
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *AuthHandler) Me(c *gin.Context) {
	userID, _, _, ok := middleware.CurrentUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var user models.User
	if err := h.DB.First(&user, "id = ?", userID).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":        user.ID.String(),
		"tenant_id": user.TenantID.String(),
		"email":     user.Email,
		"name":      user.Name,
		"role":      user.Role,
		"status":    user.Status,
	})
}
