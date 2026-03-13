package handlers

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/darkvoice1/multi-tenant-saas-platform/backend/internal/config"
	"github.com/darkvoice1/multi-tenant-saas-platform/backend/internal/middleware"
	"github.com/darkvoice1/multi-tenant-saas-platform/backend/internal/models"
	"github.com/darkvoice1/multi-tenant-saas-platform/backend/internal/observability"
	"github.com/darkvoice1/multi-tenant-saas-platform/backend/internal/storage"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/image/draw"
	"gorm.io/gorm"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"io"
	"mime"
	"net/http"
	"path"
	"strings"
	"time"
)

var errQuotaExceeded = errors.New("quota exceeded")

const (
	maxProjectNameLen  = 120
	maxProjectDescLen  = 2000
	maxTaskTitleLen    = 200
	maxCommentLen      = 2000
	maxFileSizeBytes   = 20 << 20
	maxFileNameLen     = 200
	defaultProjectName = "Untitled"
)

type CollabHandler struct {
	DB      *gorm.DB
	Config  config.Config
	Storage storage.Storage
}

func NewCollabHandler(db *gorm.DB, cfg config.Config, store storage.Storage) *CollabHandler {
	return &CollabHandler{DB: db, Config: cfg, Storage: store}
}

type projectRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type taskRequest struct {
	Title      string  `json:"title" binding:"required"`
	Status     string  `json:"status"`
	AssigneeID *string `json:"assignee_id"`
	Priority   string  `json:"priority"`
	DueAt      *string `json:"due_at"`
}

type statusRequest struct {
	Status string `json:"status" binding:"required"`
	Reason string `json:"reason"`
}

type approvalRequest struct {
	Status  string `json:"status" binding:"required"`
	Comment string `json:"comment"`
}

type commentRequest struct {
	Content string `json:"content" binding:"required"`
}

func (h *CollabHandler) ListProjects(c *gin.Context) {
	tenantID, ok := mustTenantUUID(c)
	if !ok {
		return
	}
	var projects []models.Project
	if err := h.DB.Where("tenant_id = ?", tenantID).Find(&projects).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "query failed"})
		return
	}
	c.JSON(http.StatusOK, projects)
}

func (h *CollabHandler) CreateProject(c *gin.Context) {
	tenantID, ok := mustTenantUUID(c)
	if !ok {
		return
	}
	userID, ok := mustUserUUID(c)
	if !ok {
		return
	}
	var req projectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Description = strings.TrimSpace(req.Description)
	if req.Name == "" {
		req.Name = defaultProjectName
	}
	if !validateLength(req.Name, 1, maxProjectNameLen) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project name", "code": "INVALID_NAME"})
		return
	}
	if len(req.Description) > maxProjectDescLen {
		c.JSON(http.StatusBadRequest, gin.H{"error": "description too long", "code": "INVALID_DESCRIPTION"})
		return
	}

	if err := h.enforceProjectQuota(tenantID); err != nil {
		if errors.Is(err, errQuotaExceeded) {
			c.JSON(http.StatusConflict, gin.H{"error": "project quota exceeded", "code": "PROJECT_QUOTA_EXCEEDED"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "quota check failed"})
		return
	}

	proj := models.Project{
		TenantID:    tenantID,
		Name:        req.Name,
		Description: req.Description,
		CreatedBy:   &userID,
	}
	if err := h.DB.Create(&proj).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "create failed"})
		return
	}
	c.JSON(http.StatusCreated, proj)
}

func (h *CollabHandler) GetProject(c *gin.Context) {
	tenantID, ok := mustTenantUUID(c)
	if !ok {
		return
	}
	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}
	var proj models.Project
	if err := h.DB.Where("tenant_id = ? AND id = ?", tenantID, projectID).First(&proj).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, proj)
}

func (h *CollabHandler) UpdateProject(c *gin.Context) {
	tenantID, ok := mustTenantUUID(c)
	if !ok {
		return
	}
	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}
	var req projectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	req.Name = strings.TrimSpace(req.Name)
	req.Description = strings.TrimSpace(req.Description)
	if !validateLength(req.Name, 1, maxProjectNameLen) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project name", "code": "INVALID_NAME"})
		return
	}
	if len(req.Description) > maxProjectDescLen {
		c.JSON(http.StatusBadRequest, gin.H{"error": "description too long", "code": "INVALID_DESCRIPTION"})
		return
	}

	var proj models.Project
	if err := h.DB.Where("tenant_id = ? AND id = ?", tenantID, projectID).First(&proj).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	updates := map[string]interface{}{
		"name":        req.Name,
		"description": req.Description,
	}
	if err := h.DB.Model(&proj).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "update failed"})
		return
	}
	c.JSON(http.StatusOK, proj)
}

func (h *CollabHandler) DeleteProject(c *gin.Context) {
	if !requireConfirm(c) {
		return
	}
	tenantID, ok := mustTenantUUID(c)
	if !ok {
		return
	}
	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}
	if err := h.DB.Where("tenant_id = ? AND id = ?", tenantID, projectID).Delete(&models.Project{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "delete failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *CollabHandler) ListTasksByProject(c *gin.Context) {
	tenantID, ok := mustTenantUUID(c)
	if !ok {
		return
	}
	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}
	var tasks []models.Task
	if err := h.DB.Where("tenant_id = ? AND project_id = ?", tenantID, projectID).Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "query failed"})
		return
	}
	c.JSON(http.StatusOK, tasks)
}

func (h *CollabHandler) CreateTask(c *gin.Context) {
	tenantID, ok := mustTenantUUID(c)
	if !ok {
		return
	}
	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}
	if !h.projectExists(tenantID, projectID) {
		c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
		return
	}
	var req taskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	req.Title = strings.TrimSpace(req.Title)
	if !validateLength(req.Title, 1, maxTaskTitleLen) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid title", "code": "INVALID_TITLE"})
		return
	}

	status := normalizeStatus(req.Status, "todo")
	priority := normalizePriority(req.Priority)
	var assignee *uuid.UUID
	if req.AssigneeID != nil && *req.AssigneeID != "" {
		id, err := uuid.Parse(*req.AssigneeID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid assignee_id"})
			return
		}
		assignee = &id
	}
	var dueAt *time.Time
	if req.DueAt != nil && *req.DueAt != "" {
		parsed, err := time.Parse(time.RFC3339, *req.DueAt)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid due_at"})
			return
		}
		dueAt = &parsed
	}
	task := models.Task{
		TenantID:   tenantID,
		ProjectID:  projectID,
		Title:      req.Title,
		Status:     status,
		AssigneeID: assignee,
		Priority:   priority,
		DueAt:      dueAt,
	}
	if err := h.DB.Create(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "create failed"})
		return
	}
	c.JSON(http.StatusCreated, task)
}

func (h *CollabHandler) GetTask(c *gin.Context) {
	tenantID, ok := mustTenantUUID(c)
	if !ok {
		return
	}
	taskID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}
	var task models.Task
	if err := h.DB.Where("tenant_id = ? AND id = ?", tenantID, taskID).First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, task)
}

func (h *CollabHandler) UpdateTask(c *gin.Context) {
	tenantID, ok := mustTenantUUID(c)
	if !ok {
		return
	}
	taskID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}
	var req taskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	req.Title = strings.TrimSpace(req.Title)
	if !validateLength(req.Title, 1, maxTaskTitleLen) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid title", "code": "INVALID_TITLE"})
		return
	}

	var task models.Task
	if err := h.DB.Where("tenant_id = ? AND id = ?", tenantID, taskID).First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	updates := map[string]interface{}{
		"title":    req.Title,
		"priority": normalizePriority(req.Priority),
	}
	if req.Status != "" {
		updates["status"] = normalizeStatus(req.Status, task.Status)
	}
	if req.AssigneeID != nil {
		if *req.AssigneeID == "" {
			updates["assignee_id"] = nil
		} else {
			id, err := uuid.Parse(*req.AssigneeID)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid assignee_id"})
				return
			}
			updates["assignee_id"] = id
		}
	}
	if req.DueAt != nil {
		if *req.DueAt == "" {
			updates["due_at"] = nil
		} else {
			parsed, err := time.Parse(time.RFC3339, *req.DueAt)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid due_at"})
				return
			}
			updates["due_at"] = parsed
		}
	}
	if err := h.DB.Model(&task).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "update failed"})
		return
	}
	c.JSON(http.StatusOK, task)
}

func (h *CollabHandler) DeleteTask(c *gin.Context) {
	if !requireConfirm(c) {
		return
	}
	tenantID, ok := mustTenantUUID(c)
	if !ok {
		return
	}
	taskID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}
	if err := h.DB.Where("tenant_id = ? AND id = ?", tenantID, taskID).Delete(&models.Task{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "delete failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *CollabHandler) UpdateTaskStatus(c *gin.Context) {
	tenantID, ok := mustTenantUUID(c)
	if !ok {
		return
	}
	taskID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}
	var req statusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	var task models.Task
	if err := h.DB.Where("tenant_id = ? AND id = ?", tenantID, taskID).First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	newStatus := normalizeStatus(req.Status, task.Status)
	if !isStatusTransitionAllowed(task.Status, newStatus) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status transition"})
		return
	}
	if err := h.DB.Model(&task).Update("status", newStatus).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "update failed"})
		return
	}
	c.JSON(http.StatusOK, task)
}

func (h *CollabHandler) ApproveTask(c *gin.Context) {
	tenantID, ok := mustTenantUUID(c)
	if !ok {
		return
	}
	userID, ok := mustUserUUID(c)
	if !ok {
		return
	}
	taskID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}
	var req approvalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	status := strings.ToLower(req.Status)
	if status != "approved" && status != "rejected" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status"})
		return
	}
	var task models.Task
	if err := h.DB.Where("tenant_id = ? AND id = ?", tenantID, taskID).First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	approval := models.TaskApproval{
		TenantID:   tenantID,
		TaskID:     taskID,
		ApproverID: userID,
		Status:     status,
		Comment:    req.Comment,
	}
	if err := h.DB.Create(&approval).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "approve failed"})
		return
	}
	if status == "approved" {
		observability.ObserveApprovalDuration(time.Since(task.CreatedAt))
		h.DB.Model(&task).Update("status", "done")
	} else {
		h.DB.Model(&task).Update("status", "in_progress")
	}
	if task.AssigneeID != nil {
		msg := fmt.Sprintf("Task %s approval %s", task.Title, status)
		h.createNotification(tenantID, *task.AssigneeID, "approval", msg)
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *CollabHandler) ListComments(c *gin.Context) {
	tenantID, ok := mustTenantUUID(c)
	if !ok {
		return
	}
	taskID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}
	if !h.taskExists(tenantID, taskID) {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}
	var comments []models.TaskComment
	if err := h.DB.Where("tenant_id = ? AND task_id = ?", tenantID, taskID).Find(&comments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "query failed"})
		return
	}
	c.JSON(http.StatusOK, comments)
}

func (h *CollabHandler) CreateComment(c *gin.Context) {
	tenantID, ok := mustTenantUUID(c)
	if !ok {
		return
	}
	userID, ok := mustUserUUID(c)
	if !ok {
		return
	}
	taskID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}
	if !h.taskExists(tenantID, taskID) {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}
	var req commentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	req.Content = strings.TrimSpace(req.Content)
	if !validateLength(req.Content, 1, maxCommentLen) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid content", "code": "INVALID_CONTENT"})
		return
	}

	comment := models.TaskComment{
		TenantID: tenantID,
		TaskID:   taskID,
		UserID:   userID,
		Content:  req.Content,
	}
	if err := h.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "create failed"})
		return
	}
	mentions := parseMentions(req.Content)
	if len(mentions) > 0 {
		var users []models.User
		if err := h.DB.Where("tenant_id = ? AND email IN ?", tenantID, mentions).Find(&users).Error; err == nil {
			for _, u := range users {
				msg := fmt.Sprintf("You were mentioned on task %s", taskID.String())
				h.createNotification(tenantID, u.ID, "mention", msg)
			}
		}
	}
	c.JSON(http.StatusCreated, comment)
}

func (h *CollabHandler) ListNotifications(c *gin.Context) {
	tenantID, ok := mustTenantUUID(c)
	if !ok {
		return
	}
	userID, ok := mustUserUUID(c)
	if !ok {
		return
	}
	var items []models.Notification
	if err := h.DB.Where("tenant_id = ? AND user_id = ?", tenantID, userID).Order("created_at desc").Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "query failed"})
		return
	}
	c.JSON(http.StatusOK, items)
}

func (h *CollabHandler) MarkNotificationRead(c *gin.Context) {
	tenantID, ok := mustTenantUUID(c)
	if !ok {
		return
	}
	userID, ok := mustUserUUID(c)
	if !ok {
		return
	}
	nid, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid notification id"})
		return
	}
	now := time.Now()
	if err := h.DB.Model(&models.Notification{}).
		Where("tenant_id = ? AND user_id = ? AND id = ?", tenantID, userID, nid).
		Update("read_at", &now).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "update failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *CollabHandler) ListAttachments(c *gin.Context) {
	tenantID, ok := mustTenantUUID(c)
	if !ok {
		return
	}
	taskID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}
	if !h.taskExists(tenantID, taskID) {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}
	var items []models.TaskAttachment
	if err := h.DB.Where("tenant_id = ? AND task_id = ?", tenantID, taskID).Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "query failed"})
		return
	}
	c.JSON(http.StatusOK, items)
}

func (h *CollabHandler) UploadAttachment(c *gin.Context) {
	tenantID, ok := mustTenantUUID(c)
	if !ok {
		return
	}
	userID, ok := mustUserUUID(c)
	if !ok {
		return
	}
	taskID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}
	if !h.taskExists(tenantID, taskID) {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}
	if file.Size > maxFileSizeBytes {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file too large"})
		return
	}
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "file open failed"})
		return
	}
	defer src.Close()

	data, err := io.ReadAll(src)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "file read failed"})
		return
	}
	if int64(len(data)) > maxFileSizeBytes {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file too large"})
		return
	}

	if err := h.enforceStorageQuota(tenantID, int64(len(data))); err != nil {
		if errors.Is(err, errQuotaExceeded) {
			c.JSON(http.StatusConflict, gin.H{"error": "storage quota exceeded", "code": "STORAGE_QUOTA_EXCEEDED"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "quota check failed"})
		return
	}

	contentType := file.Header.Get("Content-Type")
	if contentType == "" {
		contentType = http.DetectContentType(data)
	}
	fileName := sanitizeFileName(file.Filename)
	key := buildStorageKey(tenantID, taskID, fileName, contentType)
	if err := h.Storage.Save(c, key, bytes.NewReader(data), int64(len(data)), contentType); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "storage save failed"})
		return
	}

	var previewKey string
	if isImageContentType(contentType) {
		previewBytes, err := generateThumbnail(data)
		if err == nil && len(previewBytes) > 0 {
			previewKey = strings.TrimSuffix(key, path.Ext(key)) + "_thumb.png"
			_ = h.Storage.Save(c, previewKey, bytes.NewReader(previewBytes), int64(len(previewBytes)), "image/png")
		}
	}

	item := models.TaskAttachment{
		TenantID:    tenantID,
		TaskID:      taskID,
		UploaderID:  userID,
		FileName:    fileName,
		ContentType: contentType,
		SizeBytes:   int64(len(data)),
		Path:        key,
		PreviewPath: previewKey,
	}
	if err := h.DB.Create(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "save failed"})
		return
	}
	c.JSON(http.StatusCreated, item)
}

func (h *CollabHandler) DownloadAttachment(c *gin.Context) {
	tenantID, ok := mustTenantUUID(c)
	if !ok {
		return
	}
	if _, ok := mustUserUUID(c); !ok {
		return
	}
	attachmentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid attachment id"})
		return
	}
	var item models.TaskAttachment
	if err := h.DB.Where("tenant_id = ? AND id = ?", tenantID, attachmentID).First(&item).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	reader, err := h.Storage.Open(c, item.Path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "open failed"})
		return
	}
	defer reader.Close()
	c.Header("Content-Type", item.ContentType)
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", sanitizeFileName(item.FileName)))
	_, _ = io.Copy(c.Writer, reader)
}

func (h *CollabHandler) PreviewAttachment(c *gin.Context) {
	tenantID, ok := mustTenantUUID(c)
	if !ok {
		return
	}
	if _, ok := mustUserUUID(c); !ok {
		return
	}
	attachmentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid attachment id"})
		return
	}
	var item models.TaskAttachment
	if err := h.DB.Where("tenant_id = ? AND id = ?", tenantID, attachmentID).First(&item).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	key := item.PreviewPath
	contentType := "image/png"
	if key == "" {
		if !isImageContentType(item.ContentType) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "not previewable"})
			return
		}
		key = item.Path
		contentType = item.ContentType
	}
	reader, err := h.Storage.Open(c, key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "open failed"})
		return
	}
	defer reader.Close()
	c.Header("Content-Type", contentType)
	c.Header("Content-Disposition", "inline")
	_, _ = io.Copy(c.Writer, reader)
}

func (h *CollabHandler) projectExists(tenantID, projectID uuid.UUID) bool {
	var count int64
	h.DB.Model(&models.Project{}).Where("tenant_id = ? AND id = ?", tenantID, projectID).Count(&count)
	return count > 0
}

func (h *CollabHandler) createNotification(tenantID, userID uuid.UUID, typ, msg string) {
	n := models.Notification{TenantID: tenantID, UserID: userID, Type: typ, Message: msg}
	_ = h.DB.Create(&n).Error
}

func (h *CollabHandler) taskExists(tenantID, taskID uuid.UUID) bool {
	var count int64
	h.DB.Model(&models.Task{}).Where("tenant_id = ? AND id = ?", tenantID, taskID).Count(&count)
	return count > 0
}

func (h *CollabHandler) enforceProjectQuota(tenantID uuid.UUID) error {
	var tenant models.Tenant
	if err := h.DB.Select("id", "max_projects").Where("id = ?", tenantID).First(&tenant).Error; err != nil {
		return err
	}
	if tenant.MaxProjects <= 0 {
		return nil
	}
	var count int64
	if err := h.DB.Model(&models.Project{}).Where("tenant_id = ?", tenantID).Count(&count).Error; err != nil {
		return err
	}
	if count >= int64(tenant.MaxProjects) {
		return errQuotaExceeded
	}
	return nil
}

func (h *CollabHandler) enforceStorageQuota(tenantID uuid.UUID, incoming int64) error {
	var tenant models.Tenant
	if err := h.DB.Select("id", "max_storage_bytes").Where("id = ?", tenantID).First(&tenant).Error; err != nil {
		return err
	}
	if tenant.MaxStorageBytes <= 0 {
		return nil
	}
	var used int64
	if err := h.DB.Model(&models.TaskAttachment{}).
		Select("COALESCE(SUM(size_bytes),0)").Where("tenant_id = ?", tenantID).
		Scan(&used).Error; err != nil {
		return err
	}
	if used+incoming > tenant.MaxStorageBytes {
		return errQuotaExceeded
	}
	return nil
}

func mustTenantUUID(c *gin.Context) (uuid.UUID, bool) {
	id, ok := middleware.TenantID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing tenant id"})
		return uuid.UUID{}, false
	}
	uuidVal, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tenant id"})
		return uuid.UUID{}, false
	}
	return uuidVal, true
}

func mustUserUUID(c *gin.Context) (uuid.UUID, bool) {
	userID, _, _, ok := middleware.CurrentUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return uuid.UUID{}, false
	}
	uid, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user"})
		return uuid.UUID{}, false
	}
	return uid, true
}

func requireConfirm(c *gin.Context) bool {
	confirm := strings.ToLower(strings.TrimSpace(c.Query("confirm")))
	if confirm == "true" || confirm == "1" || confirm == "yes" {
		return true
	}
	c.JSON(http.StatusBadRequest, gin.H{"error": "confirm required", "code": "CONFIRM_REQUIRED"})
	return false
}

func normalizeStatus(status, fallback string) string {
	if status == "" {
		return fallback
	}
	s := strings.ToLower(status)
	allowed := map[string]struct{}{
		"todo":        {},
		"in_progress": {},
		"review":      {},
		"done":        {},
		"rejected":    {},
	}
	if _, ok := allowed[s]; ok {
		return s
	}
	return fallback
}

func normalizePriority(priority string) string {
	s := strings.ToLower(priority)
	if s == "" {
		return "medium"
	}
	switch s {
	case "low", "medium", "high", "urgent":
		return s
	default:
		return "medium"
	}
}

func buildStorageKey(tenantID, taskID uuid.UUID, filename, contentType string) string {
	ext := fileExtFrom(filename, contentType)
	name := uuid.New().String()
	if ext != "" {
		name += ext
	}
	return path.Join(tenantID.String(), taskID.String(), name)
}

func fileExtFrom(filename, contentType string) string {
	ext := strings.ToLower(path.Ext(filename))
	if ext != "" {
		return ext
	}
	exts, _ := mime.ExtensionsByType(contentType)
	if len(exts) > 0 {
		return exts[0]
	}
	return ""
}

func sanitizeFileName(name string) string {
	name = strings.ReplaceAll(name, "\\", "/")
	name = path.Base(name)
	name = strings.ReplaceAll(name, "\"", "")
	name = strings.ReplaceAll(name, "\n", "")
	name = strings.ReplaceAll(name, "\r", "")
	name = strings.ReplaceAll(name, "\t", "")
	name = strings.ReplaceAll(name, "..", "")
	if name == "" {
		return "attachment"
	}
	if len(name) > maxFileNameLen {
		return name[:maxFileNameLen]
	}
	return name
}

func isImageContentType(contentType string) bool {
	return strings.HasPrefix(strings.ToLower(contentType), "image/")
}

func generateThumbnail(data []byte) ([]byte, error) {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	const maxSize = 256
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	if width == 0 || height == 0 {
		return nil, fmt.Errorf("invalid image")
	}
	if width <= maxSize && height <= maxSize {
		var buf bytes.Buffer
		if err := png.Encode(&buf, img); err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	}
	ratio := float64(maxSize) / float64(width)
	if height > width {
		ratio = float64(maxSize) / float64(height)
	}
	newW := int(float64(width) * ratio)
	newH := int(float64(height) * ratio)
	dst := image.NewRGBA(image.Rect(0, 0, newW, newH))
	draw.CatmullRom.Scale(dst, dst.Bounds(), img, bounds, draw.Over, nil)
	var buf bytes.Buffer
	if err := png.Encode(&buf, dst); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func isStatusTransitionAllowed(from, to string) bool {
	if from == to {
		return true
	}
	allowed := map[string][]string{
		"todo":        {"in_progress"},
		"in_progress": {"review", "todo"},
		"review":      {"done", "in_progress", "rejected"},
		"rejected":    {"in_progress"},
		"done":        {},
	}
	for _, n := range allowed[from] {
		if n == to {
			return true
		}
	}
	return false
}

func parseMentions(content string) []string {
	seen := map[string]struct{}{}
	var mentions []string
	parts := strings.Fields(content)
	for _, p := range parts {
		if strings.HasPrefix(p, "@") {
			email := strings.Trim(p, "@,.;:")
			if strings.Contains(email, "@") && strings.Contains(email, ".") {
				if _, ok := seen[email]; !ok {
					seen[email] = struct{}{}
					mentions = append(mentions, email)
				}
			}
		}
	}
	return mentions
}

func validateLength(value string, min, max int) bool {
	l := len(value)
	return l >= min && l <= max
}
