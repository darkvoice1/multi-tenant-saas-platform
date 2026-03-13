package middleware

import (
	"strconv"
	"sync"
	"time"

	"github.com/darkvoice1/multi-tenant-saas-platform/backend/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type tenantRateState struct {
	Count          int
	ResetAt        time.Time
	Limit          int
	LimitFetchedAt time.Time
}

var (
	rateMu       sync.Mutex
	rateByTenant = map[string]*tenantRateState{}
)

const defaultMaxRequestsPerMinute = 600

func RateLimit(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID, ok := TenantID(c)
		if !ok || tenantID == "" {
			c.Next()
			return
		}

		now := time.Now()
		state := getTenantRateState(db, tenantID, now)

		if state.Limit <= 0 {
			c.Next()
			return
		}

		if now.After(state.ResetAt) {
			state.Count = 0
			state.ResetAt = now.Add(time.Minute)
		}

		state.Count++
		if state.Count > state.Limit {
			retryAfter := int(state.ResetAt.Sub(now).Seconds())
			if retryAfter < 1 {
				retryAfter = 1
			}
			c.Header("Retry-After", strconv.Itoa(retryAfter))
			c.JSON(429, gin.H{
				"error": "rate limit exceeded",
				"code":  "RATE_LIMITED",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func getTenantRateState(db *gorm.DB, tenantID string, now time.Time) *tenantRateState {
	rateMu.Lock()
	defer rateMu.Unlock()

	state, ok := rateByTenant[tenantID]
	if !ok {
		state = &tenantRateState{
			ResetAt: now.Add(time.Minute),
			Limit:   defaultMaxRequestsPerMinute,
		}
		rateByTenant[tenantID] = state
	}

	if state.LimitFetchedAt.IsZero() || now.Sub(state.LimitFetchedAt) > time.Minute {
		var tenant models.Tenant
		if err := db.Select("max_requests_per_minute").Where("id = ?", tenantID).First(&tenant).Error; err == nil {
			if tenant.MaxRequestsPerMinute > 0 {
				state.Limit = tenant.MaxRequestsPerMinute
			} else {
				state.Limit = defaultMaxRequestsPerMinute
			}
		}
		state.LimitFetchedAt = now
	}

	return state
}
