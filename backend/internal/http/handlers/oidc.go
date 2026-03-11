package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/darkvoice1/multi-tenant-saas-platform/backend/internal/config"
	"github.com/darkvoice1/multi-tenant-saas-platform/backend/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OIDCHandler struct {
	DB     *gorm.DB
	Config config.Config
	store  *oidcStore
}

type oidcStore struct {
	mu    sync.Mutex
	codes map[string]oidcCode
}

type oidcCode struct {
	TenantID  uuid.UUID
	UserID    uuid.UUID
	Role      string
	State     string
	ExpiresAt time.Time
}

func newOIDCStore() *oidcStore {
	return &oidcStore{codes: map[string]oidcCode{}}
}

func (s *oidcStore) put(code string, data oidcCode) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.codes[code] = data
}

func (s *oidcStore) pop(code string) (oidcCode, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	data, ok := s.codes[code]
	if ok {
		delete(s.codes, code)
	}
	return data, ok
}

func NewOIDCHandler(db *gorm.DB, cfg config.Config) *OIDCHandler {
	return &OIDCHandler{DB: db, Config: cfg, store: newOIDCStore()}
}

func (h *OIDCHandler) MockAuthorize(c *gin.Context) {
	tenantIDStr := c.Query("tenant_id")
	email := c.Query("email")
	state := c.Query("state")
	redirectURI := c.Query("redirect_uri")

	if tenantIDStr == "" || email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tenant_id and email are required"})
		return
	}

	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tenant_id"})
		return
	}

	var user models.User
	if err := h.DB.Where("tenant_id = ? AND email = ?", tenantID, email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	code, err := generateCode()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "code generation failed"})
		return
	}

	h.store.put(code, oidcCode{
		TenantID:  user.TenantID,
		UserID:    user.ID,
		Role:      user.Role,
		State:     state,
		ExpiresAt: time.Now().Add(5 * time.Minute),
	})

	if redirectURI != "" {
		u, err := url.Parse(redirectURI)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid redirect_uri"})
			return
		}
		q := u.Query()
		q.Set("code", code)
		if state != "" {
			q.Set("state", state)
		}
		u.RawQuery = q.Encode()
		c.Redirect(http.StatusFound, u.String())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":  code,
		"state": state,
	})
}

type oidcCallbackRequest struct {
	Code  string `json:"code"`
	State string `json:"state"`
}

func (h *OIDCHandler) Callback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")
	if code == "" {
		var req oidcCallbackRequest
		if err := c.ShouldBindJSON(&req); err == nil {
			code = req.Code
			state = req.State
		}
	}

	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing code"})
		return
	}

	data, ok := h.store.pop(code)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid code"})
		return
	}
	if data.ExpiresAt.Before(time.Now()) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "code expired"})
		return
	}
	if data.State != "" && state != "" && !strings.EqualFold(data.State, state) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "state mismatch"})
		return
	}

	var user models.User
	if err := h.DB.First(&user, "id = ?", data.UserID).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user"})
		return
	}
	if user.Status != "active" {
		c.JSON(http.StatusForbidden, gin.H{"error": "user disabled"})
		return
	}

	accessToken, refreshToken, expiresIn, err := issueTokens(h.DB, h.Config, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token error"})
		return
	}

	now := time.Now()
	h.DB.Model(&user).Updates(map[string]interface{}{
		"last_login_at": &now,
	})

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

func generateCode() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}
