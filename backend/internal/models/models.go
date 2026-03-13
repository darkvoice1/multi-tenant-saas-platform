package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Base includes common fields for all models.
type Base struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CreatedAt time.Time      `json:"created_at" gorm:"not null;default:now()"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"not null;default:now()"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type Tenant struct {
	Base
	Name                 string `json:"name" gorm:"not null"`
	Slug                 string `json:"slug" gorm:"not null;uniqueIndex"`
	Status               string `json:"status" gorm:"not null;default:active"`
	MaxProjects          int    `json:"max_projects" gorm:"not null;default:100"`
	MaxMembers           int    `json:"max_members" gorm:"not null;default:100"`
	MaxStorageBytes      int64  `json:"max_storage_bytes" gorm:"not null;default:1073741824"`
	MaxRequestsPerMinute int    `json:"max_requests_per_minute" gorm:"not null;default:600"`
}

type Org struct {
	Base
	TenantID uuid.UUID `json:"tenant_id" gorm:"type:uuid;not null;index"`
	Name     string    `json:"name" gorm:"not null"`
}

type User struct {
	Base
	TenantID     uuid.UUID  `json:"tenant_id" gorm:"type:uuid;not null;index"`
	OrgID        *uuid.UUID `json:"org_id,omitempty" gorm:"type:uuid;index"`
	Email        string     `json:"email" gorm:"not null"`
	Name         string     `json:"name" gorm:"not null"`
	Role         string     `json:"role" gorm:"not null"`
	Status       string     `json:"status" gorm:"not null;default:active"`
	PasswordHash string     `json:"-" gorm:"type:text"`
	LastLoginAt  *time.Time `json:"last_login_at,omitempty"`
}

type Project struct {
	Base
	TenantID    uuid.UUID  `json:"tenant_id" gorm:"type:uuid;not null;index"`
	OrgID       *uuid.UUID `json:"org_id,omitempty" gorm:"type:uuid;index"`
	Name        string     `json:"name" gorm:"not null"`
	Description string     `json:"description"`
	CreatedBy   *uuid.UUID `json:"created_by,omitempty" gorm:"type:uuid"`
}

type Task struct {
	Base
	TenantID   uuid.UUID  `json:"tenant_id" gorm:"type:uuid;not null;index"`
	ProjectID  uuid.UUID  `json:"project_id" gorm:"type:uuid;not null;index"`
	Title      string     `json:"title" gorm:"not null"`
	Status     string     `json:"status" gorm:"not null;default:todo"`
	AssigneeID *uuid.UUID `json:"assignee_id,omitempty" gorm:"type:uuid"`
	Priority   string     `json:"priority" gorm:"not null;default:medium"`
	DueAt      *time.Time `json:"due_at,omitempty"`
}

type RefreshToken struct {
	Base
	TenantID  uuid.UUID  `json:"tenant_id" gorm:"type:uuid;not null;index"`
	UserID    uuid.UUID  `json:"user_id" gorm:"type:uuid;not null;index"`
	TokenHash string     `json:"-" gorm:"type:text;not null;uniqueIndex"`
	ExpiresAt time.Time  `json:"expires_at" gorm:"not null"`
	RevokedAt *time.Time `json:"revoked_at,omitempty"`
}

type TaskComment struct {
	Base
	TenantID uuid.UUID `json:"tenant_id" gorm:"type:uuid;not null;index"`
	TaskID   uuid.UUID `json:"task_id" gorm:"type:uuid;not null;index"`
	UserID   uuid.UUID `json:"user_id" gorm:"type:uuid;not null;index"`
	Content  string    `json:"content" gorm:"type:text;not null"`
}

type TaskApproval struct {
	Base
	TenantID   uuid.UUID `json:"tenant_id" gorm:"type:uuid;not null;index"`
	TaskID     uuid.UUID `json:"task_id" gorm:"type:uuid;not null;index"`
	ApproverID uuid.UUID `json:"approver_id" gorm:"type:uuid;not null;index"`
	Status     string    `json:"status" gorm:"not null"` // approved / rejected
	Comment    string    `json:"comment" gorm:"type:text"`
}

type TaskAttachment struct {
	Base
	TenantID    uuid.UUID `json:"tenant_id" gorm:"type:uuid;not null;index"`
	TaskID      uuid.UUID `json:"task_id" gorm:"type:uuid;not null;index"`
	UploaderID  uuid.UUID `json:"uploader_id" gorm:"type:uuid;not null;index"`
	FileName    string    `json:"file_name" gorm:"not null"`
	ContentType string    `json:"content_type" gorm:"not null"`
	SizeBytes   int64     `json:"size_bytes" gorm:"not null"`
	Path        string    `json:"-" gorm:"type:text;not null"`
	PreviewPath string    `json:"-" gorm:"type:text"`
}

type Notification struct {
	Base
	TenantID uuid.UUID  `json:"tenant_id" gorm:"type:uuid;not null;index"`
	UserID   uuid.UUID  `json:"user_id" gorm:"type:uuid;not null;index"`
	Type     string     `json:"type" gorm:"not null"`
	Message  string     `json:"message" gorm:"type:text;not null"`
	ReadAt   *time.Time `json:"read_at,omitempty"`
}

type AuditLog struct {
	Base
	TenantID   uuid.UUID  `json:"tenant_id" gorm:"type:uuid;not null;index"`
	UserID     uuid.UUID  `json:"user_id" gorm:"type:uuid;not null;index"`
	Action     string     `json:"action" gorm:"not null"`
	Resource   string     `json:"resource" gorm:"not null"`
	ResourceID *uuid.UUID `json:"resource_id,omitempty" gorm:"type:uuid"`
	Method     string     `json:"method" gorm:"not null"`
	Path       string     `json:"path" gorm:"not null"`
	StatusCode int        `json:"status_code" gorm:"not null"`
	IP         string     `json:"ip" gorm:"type:text"`
	UserAgent  string     `json:"user_agent" gorm:"type:text"`
}
