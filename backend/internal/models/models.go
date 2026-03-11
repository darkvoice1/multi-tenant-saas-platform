package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Base includes common fields for all models.
type Base struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CreatedAt time.Time      `gorm:"not null;default:now()"`
	UpdatedAt time.Time      `gorm:"not null;default:now()"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Tenant struct {
	Base
	Name   string `gorm:"not null"`
	Slug   string `gorm:"not null;uniqueIndex"`
	Status string `gorm:"not null;default:active"`
}

type Org struct {
	Base
	TenantID uuid.UUID `gorm:"type:uuid;not null;index"`
	Name     string    `gorm:"not null"`
}

type User struct {
	Base
	TenantID     uuid.UUID  `gorm:"type:uuid;not null;index"`
	OrgID        *uuid.UUID `gorm:"type:uuid;index"`
	Email        string     `gorm:"not null"`
	Name         string     `gorm:"not null"`
	Role         string     `gorm:"not null"`
	Status       string     `gorm:"not null;default:active"`
	PasswordHash string     `gorm:"type:text"`
	LastLoginAt  *time.Time
}

type Project struct {
	Base
	TenantID    uuid.UUID  `gorm:"type:uuid;not null;index"`
	OrgID       *uuid.UUID `gorm:"type:uuid;index"`
	Name        string     `gorm:"not null"`
	Description string
	CreatedBy   *uuid.UUID `gorm:"type:uuid"`
}

type Task struct {
	Base
	TenantID   uuid.UUID  `gorm:"type:uuid;not null;index"`
	ProjectID  uuid.UUID  `gorm:"type:uuid;not null;index"`
	Title      string     `gorm:"not null"`
	Status     string     `gorm:"not null;default:todo"`
	AssigneeID *uuid.UUID `gorm:"type:uuid"`
}

type RefreshToken struct {
	Base
	TenantID  uuid.UUID `gorm:"type:uuid;not null;index"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index"`
	TokenHash string    `gorm:"type:text;not null;uniqueIndex"`
	ExpiresAt time.Time `gorm:"not null"`
	RevokedAt *time.Time
}
