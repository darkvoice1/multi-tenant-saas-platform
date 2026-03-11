package http

import (
	"github.com/darkvoice1/multi-tenant-saas-platform/backend/internal/http/handlers"
	"github.com/darkvoice1/multi-tenant-saas-platform/backend/internal/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewRouter(db *gorm.DB) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	router.GET("/healthz", handlers.Health)

	api := router.Group("/api")
	api.Use(middleware.TenantMiddleware(true))
	api.GET("/tenant/echo", handlers.TenantEcho)

	_ = db

	return router
}
