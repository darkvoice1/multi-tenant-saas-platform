package http

import (
	"github.com/darkvoice1/multi-tenant-saas-platform/backend/internal/auth"
	"github.com/darkvoice1/multi-tenant-saas-platform/backend/internal/config"
	"github.com/darkvoice1/multi-tenant-saas-platform/backend/internal/http/handlers"
	"github.com/darkvoice1/multi-tenant-saas-platform/backend/internal/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewRouter(db *gorm.DB, cfg config.Config) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery(), middleware.CORSMiddleware())

	router.GET("/healthz", handlers.Health)

	authHandler := handlers.NewAuthHandler(db, cfg)
	oidcHandler := handlers.NewOIDCHandler(db, cfg)
	router.POST("/auth/login", authHandler.Login)
	router.POST("/auth/refresh", authHandler.Refresh)
	router.POST("/auth/logout", authHandler.Logout)
	router.GET("/auth/me", middleware.AuthMiddleware(cfg.JWTSecret), authHandler.Me)
	router.GET("/auth/oidc/mock/authorize", oidcHandler.MockAuthorize)
	router.GET("/auth/oidc/callback", oidcHandler.Callback)
	router.POST("/auth/oidc/callback", oidcHandler.Callback)

	api := router.Group("/api")
	api.Use(middleware.TenantMiddleware(true))
	api.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	api.GET("/tenant/echo", handlers.TenantEcho)
	api.GET("/admin/ping", middleware.RequireRoles(auth.RoleAdmin), handlers.AdminPing)

	return router
}
