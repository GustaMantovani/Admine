package api

import (
	"log/slog"
	"strings"

	"github.com/GustaMantovani/Admine/server_handler/internal"
	"github.com/GustaMantovani/Admine/server_handler/internal/api/handlers"
	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
)

// SetupRoutes configures all API routes
func SetupRoutes() *gin.Engine {
	// Set Gin mode BEFORE creating router
	logLevel := strings.ToUpper(internal.Get().Config.App.LogLevel)

	if logLevel == "DEBUG" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create Gin router AFTER setting mode
	router := gin.New()

	router.Use(sloggin.New(slog.Default()))

	// Create handlers
	serverHandler := handlers.NewApiHandler()

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
