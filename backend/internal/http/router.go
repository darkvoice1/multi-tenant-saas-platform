package http

import (
	"github.com/darkvoice1/multi-tenant-saas-platform/backend/internal/auth"
	"github.com/darkvoice1/multi-tenant-saas-platform/backend/internal/config"
	"github.com/darkvoice1/multi-tenant-saas-platform/backend/internal/http/handlers"
	"github.com/darkvoice1/multi-tenant-saas-platform/backend/internal/middleware"
	"github.com/darkvoice1/multi-tenant-saas-platform/backend/internal/storage"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewRouter(db *gorm.DB, cfg config.Config, store storage.Storage) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery(), middleware.CORSMiddleware())

	router.GET("/healthz", handlers.Health)

	authHandler := handlers.NewAuthHandler(db, cfg)
	oidcHandler := handlers.NewOIDCHandler(db, cfg)
	collabHandler := handlers.NewCollabHandler(db, cfg, store)
	dashboardHandler := handlers.NewDashboardHandler(db, cfg)
	router.POST("/auth/bootstrap", authHandler.Bootstrap)
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
	api.GET("/dashboard", middleware.RequirePermission(auth.PermProjectRead), dashboardHandler.Get)

	projects := api.Group("/projects")
	projects.GET("", middleware.RequirePermission(auth.PermProjectRead), collabHandler.ListProjects)
	projects.POST("", middleware.RequirePermission(auth.PermProjectWrite), collabHandler.CreateProject)
	projects.GET("/:id", middleware.RequirePermission(auth.PermProjectRead), collabHandler.GetProject)
	projects.PUT("/:id", middleware.RequirePermission(auth.PermProjectWrite), collabHandler.UpdateProject)
	projects.DELETE("/:id", middleware.RequirePermission(auth.PermProjectWrite), collabHandler.DeleteProject)
	projects.GET("/:id/tasks", middleware.RequirePermission(auth.PermTaskRead), collabHandler.ListTasksByProject)
	projects.POST("/:id/tasks", middleware.RequirePermission(auth.PermTaskWrite), collabHandler.CreateTask)

	api.GET("/tasks/:id", middleware.RequirePermission(auth.PermTaskRead), collabHandler.GetTask)
	api.PUT("/tasks/:id", middleware.RequirePermission(auth.PermTaskWrite), collabHandler.UpdateTask)
	api.DELETE("/tasks/:id", middleware.RequirePermission(auth.PermTaskWrite), collabHandler.DeleteTask)
	api.POST("/tasks/:id/status", middleware.RequirePermission(auth.PermTaskWrite), collabHandler.UpdateTaskStatus)
	api.POST("/tasks/:id/approve", middleware.RequireRoles(auth.RoleAdmin, auth.RoleManager), collabHandler.ApproveTask)
	api.GET("/tasks/:id/comments", middleware.RequirePermission(auth.PermTaskRead), collabHandler.ListComments)
	api.POST("/tasks/:id/comments", middleware.RequirePermission(auth.PermTaskWrite), collabHandler.CreateComment)
	api.GET("/tasks/:id/attachments", middleware.RequirePermission(auth.PermTaskRead), collabHandler.ListAttachments)
	api.POST("/tasks/:id/attachments", middleware.RequirePermission(auth.PermTaskWrite), collabHandler.UploadAttachment)
	api.GET("/attachments/:id/download", middleware.RequirePermission(auth.PermTaskRead), collabHandler.DownloadAttachment)
	api.GET("/attachments/:id/preview", middleware.RequirePermission(auth.PermTaskRead), collabHandler.PreviewAttachment)

	api.GET("/notifications", collabHandler.ListNotifications)
	api.POST("/notifications/:id/read", collabHandler.MarkNotificationRead)

	return router
}
