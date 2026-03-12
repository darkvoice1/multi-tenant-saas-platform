package main

import (
	"fmt"
	"log"
	"os"

	"github.com/darkvoice1/multi-tenant-saas-platform/backend/internal/auth"
	"github.com/darkvoice1/multi-tenant-saas-platform/backend/internal/config"
	"github.com/darkvoice1/multi-tenant-saas-platform/backend/internal/db"
	"github.com/darkvoice1/multi-tenant-saas-platform/backend/internal/models"
	"gorm.io/gorm"
)

func main() {
	if _, err := os.Stat(".env"); err == nil {
		_ = config.LoadDotEnv(".env")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	database, err := db.Connect(cfg)
	if err != nil {
		log.Fatalf("db error: %v", err)
	}

	tenantSlug := getEnv("SEED_TENANT_SLUG", "demo")
	tenantName := getEnv("SEED_TENANT_NAME", "Demo Tenant")
	adminEmail := getEnv("SEED_ADMIN_EMAIL", "admin@example.com")
	adminPass := getEnv("SEED_ADMIN_PASSWORD", "Admin123!")
	adminName := getEnv("SEED_ADMIN_NAME", "Admin")

	var tenant models.Tenant
	if err := database.Where("slug = ?", tenantSlug).First(&tenant).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			tenant = models.Tenant{Name: tenantName, Slug: tenantSlug, Status: "active"}
			if err := database.Create(&tenant).Error; err != nil {
				log.Fatalf("create tenant: %v", err)
			}
		} else {
			log.Fatalf("query tenant: %v", err)
		}
	}

	var user models.User
	if err := database.Where("tenant_id = ? AND email = ?", tenant.ID, adminEmail).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			hash, err := auth.HashPassword(adminPass)
			if err != nil {
				log.Fatalf("hash password: %v", err)
			}
			user = models.User{
				TenantID:     tenant.ID,
				Email:        adminEmail,
				Name:         adminName,
				Role:         auth.RoleAdmin,
				Status:       "active",
				PasswordHash: hash,
			}
			if err := database.Create(&user).Error; err != nil {
				log.Fatalf("create user: %v", err)
			}
		} else {
			log.Fatalf("query user: %v", err)
		}
	}

	fmt.Printf("seeded tenant_id=%s email=%s password=%s\n", tenant.ID.String(), adminEmail, adminPass)
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
