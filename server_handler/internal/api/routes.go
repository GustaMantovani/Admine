package api

import (
	"admine.com/server_handler/internal/api/handlers"
	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all API routes
func SetupRoutes() *gin.Engine {
	// Create Gin router
	router := gin.Default()

	// Create handlers
	serverHandler := handlers.NewServerHandler()

	// API routes
	api := router.Group("/api/v1")
	{
		api.GET("/info", serverHandler.GetInfo)
		api.GET("/status", serverHandler.GetStatus)
		api.POST("/command", serverHandler.PostCommand)
	}

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "server_handler",
		})
	})

	return router
}
