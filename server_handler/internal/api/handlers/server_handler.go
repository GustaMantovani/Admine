package handlers

import (
	"log/slog"
	"net/http"

	"github.com/GustaMantovani/Admine/server_handler/internal"
	"github.com/GustaMantovani/Admine/server_handler/internal/api/models"
	"github.com/gin-gonic/gin"
)

// ServerHandler handles server-related API endpoints
type ServerHandler struct {
}

// NewApiHandler creates a new ServerHandler instance
func NewApiHandler() *ServerHandler {
	return &ServerHandler{}
}

// GetInfo handles GET /info endpoint
func (h *ServerHandler) GetInfo(c *gin.Context) {
	slog.Info("GET /info endpoint called")

	ctx := internal.Get()
	if ctx.MinecraftServer == nil {
		slog.Error("MinecraftServer is not initialized")
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Minecraft server not initialized"))
		return
	}

	// Get server info through the MinecraftServer interface
	serverInfo, err := (*ctx.MinecraftServer).Info(*ctx.MainCtx)
	if err != nil {
		slog.Error("Failed to get server info", "error", err.Error())
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to get server info: "+err.Error()))
		return
	}

	slog.Info("Successfully retrieved server info")
	c.JSON(http.StatusOK, serverInfo)
}

// GetStatus handles GET /status endpoint
func (h *ServerHandler) GetStatus(c *gin.Context) {
	slog.Info("GET /status endpoint called")

	ctx := internal.Get()
	if ctx.MinecraftServer == nil {
		slog.Error("MinecraftServer is not initialized")
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Minecraft server not initialized"))
		return
	}

	// Get server status through the MinecraftServer interface
	serverStatus, err := (*ctx.MinecraftServer).Status(*ctx.MainCtx)
	if err != nil {
		slog.Error("Failed to get server status", "error", err.Error())
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to get server status: "+err.Error()))
		return
	}

	slog.Info("Successfully retrieved server status")
	c.JSON(http.StatusOK, serverStatus)
}

// PostCommand handles POST /command endpoint
func (h *ServerHandler) PostCommand(c *gin.Context) {
	slog.Info("POST /command endpoint called")

	var command models.Command
	if err := c.ShouldBindJSON(&command); err != nil {
		slog.Error("Invalid command request", "error", err.Error())
		c.JSON(http.StatusBadRequest, models.NewErrorResponse("Invalid request: "+err.Error()))
		return
	}

	ctx := internal.Get()
	if ctx.MinecraftServer == nil {
		slog.Error("MinecraftServer is not initialized")
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Minecraft server not initialized"))
		return
	}

	// Execute command through the MinecraftServer interface
	result, err := (*ctx.MinecraftServer).ExecuteCommand(*ctx.MainCtx, command.Command)
	if err != nil {
		slog.Error("Failed to execute command", "command", command.Command, "error", err.Error())
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to execute command: "+err.Error()))
		return
	}

	slog.Info("Successfully executed command", "command", command.Command)
	c.JSON(http.StatusOK, result)
}
