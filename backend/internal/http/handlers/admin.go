package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AdminPing(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "admin_ok",
	})
}
