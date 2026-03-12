package handlers

import (
	"net/http"
	"time"

	"github.com/darkvoice1/multi-tenant-saas-platform/backend/internal/config"
	"github.com/darkvoice1/multi-tenant-saas-platform/backend/internal/middleware"
	"github.com/darkvoice1/multi-tenant-saas-platform/backend/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DashboardHandler struct {
	DB     *gorm.DB
	Config config.Config
}

func NewDashboardHandler(db *gorm.DB, cfg config.Config) *DashboardHandler {
	return &DashboardHandler{DB: db, Config: cfg}
}

type statusCountRow struct {
	Status string `json:"status"`
	Count  int64  `json:"count"`
}

type recentCommentItem struct {
	ID        uuid.UUID `json:"id"`
	TaskID    uuid.UUID `json:"task_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UserName  string    `json:"user_name"`
	UserEmail string    `json:"user_email"`
}

type recentTaskItem struct {
	ID        uuid.UUID  `json:"id"`
	ProjectID uuid.UUID  `json:"project_id"`
	Title     string     `json:"title"`
	Status    string     `json:"status"`
	Priority  string     `json:"priority"`
	DueAt     *time.Time `json:"due_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type recentNotificationItem struct {
	ID        uuid.UUID  `json:"id"`
	Type      string     `json:"type"`
	Message   string     `json:"message"`
	ReadAt    *time.Time `json:"read_at"`
	CreatedAt time.Time  `json:"created_at"`
}

func (h *DashboardHandler) Get(c *gin.Context) {
	tenantIDStr, ok := middleware.TenantID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing tenant id"})
		return
	}
	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tenant id"})
		return
	}

	var tenant models.Tenant
	if err := h.DB.Where("id = ?", tenantID).First(&tenant).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "tenant not found"})
		return
	}

	var currentUserID uuid.UUID
	var hasUser bool
	if uidStr, _, _, ok := middleware.CurrentUser(c); ok {
		if uid, err := uuid.Parse(uidStr); err == nil {
			currentUserID = uid
			hasUser = true
		}
	}

	var projectCount int64
	_ = h.DB.Model(&models.Project{}).Where("tenant_id = ?", tenantID).Count(&projectCount).Error

	var taskCount int64
	_ = h.DB.Model(&models.Task{}).Where("tenant_id = ?", tenantID).Count(&taskCount).Error

	statusCounts := make([]statusCountRow, 0)
	_ = h.DB.Model(&models.Task{}).
		Select("status, count(*) as count").
		Where("tenant_id = ?", tenantID).
		Group("status").
		Scan(&statusCounts).Error

	var unreadNotifications int64
	if hasUser {
		_ = h.DB.Model(&models.Notification{}).
			Where("tenant_id = ? AND user_id = ? AND read_at IS NULL", tenantID, currentUserID).
			Count(&unreadNotifications).Error
	}

	var pendingReview int64
	_ = h.DB.Model(&models.Task{}).
		Where("tenant_id = ? AND status = ?", tenantID, "review").
		Count(&pendingReview).Error

	var myOpenTasks int64
	if hasUser {
		_ = h.DB.Model(&models.Task{}).
			Where("tenant_id = ? AND assignee_id = ? AND status <> ?", tenantID, currentUserID, "done").
			Count(&myOpenTasks).Error
	}

	dueSoon := make([]recentTaskItem, 0)
	deadline := time.Now().Add(7 * 24 * time.Hour)
	_ = h.DB.Model(&models.Task{}).
		Select("id, project_id, title, status, priority, due_at, updated_at").
		Where("tenant_id = ? AND due_at IS NOT NULL AND due_at <= ? AND status <> ?", tenantID, deadline, "done").
		Order("due_at asc").
		Limit(5).
		Scan(&dueSoon).Error

	myDueSoon := make([]recentTaskItem, 0)
	if hasUser {
		_ = h.DB.Model(&models.Task{}).
			Select("id, project_id, title, status, priority, due_at, updated_at").
			Where("tenant_id = ? AND assignee_id = ? AND due_at IS NOT NULL AND due_at <= ? AND status <> ?", tenantID, currentUserID, deadline, "done").
			Order("due_at asc").
			Limit(5).
			Scan(&myDueSoon).Error
	}

	recentTasks := make([]recentTaskItem, 0)
	_ = h.DB.Model(&models.Task{}).
		Select("id, project_id, title, status, priority, due_at, updated_at").
		Where("tenant_id = ?", tenantID).
		Order("updated_at desc").
		Limit(8).
		Scan(&recentTasks).Error

	pendingReviewTasks := make([]recentTaskItem, 0)
	_ = h.DB.Model(&models.Task{}).
		Select("id, project_id, title, status, priority, due_at, updated_at").
		Where("tenant_id = ? AND status = ?", tenantID, "review").
		Order("updated_at desc").
		Limit(5).
		Scan(&pendingReviewTasks).Error

	recentComments := make([]recentCommentItem, 0)
	_ = h.DB.Table("task_comments tc").
		Select("tc.id, tc.task_id, tc.content, tc.created_at, u.name as user_name, u.email as user_email").
		Joins("join users u on u.id = tc.user_id").
		Where("tc.tenant_id = ? AND tc.deleted_at IS NULL", tenantID).
		Order("tc.created_at desc").
		Limit(8).
		Scan(&recentComments).Error

	recentNotifications := make([]recentNotificationItem, 0)
	if hasUser {
		_ = h.DB.Model(&models.Notification{}).
			Select("id, type, message, read_at, created_at").
			Where("tenant_id = ? AND user_id = ?", tenantID, currentUserID).
			Order("created_at desc").
			Limit(8).
			Scan(&recentNotifications).Error
	}

	c.JSON(http.StatusOK, gin.H{
		"tenant": gin.H{
			"id":     tenant.ID.String(),
			"name":   tenant.Name,
			"slug":   tenant.Slug,
			"status": tenant.Status,
		},
		"metrics": gin.H{
			"project_count":             projectCount,
			"task_count":                taskCount,
			"task_status_counts":        statusCounts,
			"pending_review_task_count": pendingReview,
			"unread_notification_count": unreadNotifications,
			"my_open_task_count":        myOpenTasks,
		},
		"lists": gin.H{
			"due_soon_tasks":        dueSoon,
			"my_due_soon_tasks":     myDueSoon,
			"recent_tasks":          recentTasks,
			"pending_review_tasks":  pendingReviewTasks,
			"recent_comments":       recentComments,
			"recent_notifications":  recentNotifications,
			"storage_backend":       h.Config.StorageBackend,
			"storage_s3_bucket":     h.Config.S3Bucket,
		},
	})
}
