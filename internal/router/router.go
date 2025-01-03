package router

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitRoutes(r *gin.Engine, db *gorm.DB) {
	// middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Healtcheck endpoint
	r.GET("/health-check", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "ok",
		})
	})
}
