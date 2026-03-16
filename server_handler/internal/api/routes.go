package api

import (
	"bytes"
	"context"
	"log/slog"
	"strings"

	"github.com/GustaMantovani/Admine/server_handler/internal/api/handlers"
	"github.com/GustaMantovani/Admine/server_handler/internal/config"
	"github.com/GustaMantovani/Admine/server_handler/internal/pubsub"
	"github.com/GustaMantovani/Admine/server_handler/internal/server"
	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
)

// bodyLogWriter wraps gin.ResponseWriter to capture the response body for debug logging.
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// responseBodyLogger is a middleware that logs the full JSON response body at DEBUG level.
func responseBodyLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		blw := &bodyLogWriter{body: &bytes.Buffer{}, ResponseWriter: c.Writer}
		c.Writer = blw
		c.Next()
		slog.Debug("HTTP response",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"body", blw.body.String(),
		)
	}
}

// SetupRouter configures the Gin router with all API routes
func SetupRouter(
	srv server.MinecraftServer,
	ps pubsub.PubSubService,
	origin string,
	serverChannel string,
	logLevel string,
	cfg config.MinecraftServerConfig,
	mainCtx context.Context,
) *gin.Engine {
	if strings.ToUpper(logLevel) == "DEBUG" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(sloggin.New(slog.Default()))
	if strings.ToUpper(logLevel) == "DEBUG" {
		router.Use(responseBodyLogger())
	}

	serverHandler := handlers.NewServerHandler(srv)
	modHandler := handlers.NewModHandler(srv, ps, origin, serverChannel, cfg, mainCtx)

	api := router.Group("/api/v1")
	{
		api.GET("/info", serverHandler.GetInfo)
		api.GET("/status", serverHandler.GetStatus)
		api.GET("/logs", serverHandler.GetLogs)
		api.POST("/command", serverHandler.PostCommand)
		api.GET("/resources", serverHandler.GetResourceUsage)
		api.POST("/mods", modHandler.PostInstallMod)
		api.GET("/mods", modHandler.GetListMods)
		api.DELETE("/mods/:filename", modHandler.DeleteRemoveMod)
	}

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "server_handler",
		})
	})

	return router
}
