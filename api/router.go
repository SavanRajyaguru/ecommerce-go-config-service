package api

import (
	"github.com/gin-gonic/gin"
)

// SetupRouter initializes the Gin router
func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "UP",
		})
	})

	return r
}
