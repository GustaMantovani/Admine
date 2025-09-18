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
	_, err := (*ctx.MinecraftServer).Info()
	if err != nil {
		slog.Error("Failed to get server info", "error", err.Error())
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to get server info: "+err.Error()))
		return
	}

	// For now, return static info. In a real implementation,
	// you would parse the info string or get these values from the server
	serverInfo := models.NewServerInfo(
		"1.20.1",  // MinecraftVersion - would be parsed from server
		"17.0.2",  // JavaVersion - would be detected
		"Fabric",  // ModEngine - from config or detection
		"default", // Seed - from server properties
		20,        // MaxPlayers - from server properties
	)

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
	statusStr, err := (*ctx.MinecraftServer).Status()
	if err != nil {
		slog.Error("Failed to get server status", "error", err.Error())
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to get server status: "+err.Error()))
		return
	}

	// For now, return static status. In a real implementation,
	// you would parse the status string and determine actual values
	var status models.ServerStatusEnum
	var health models.HealthStatus

	// Simple status mapping based on response
	if statusStr == "running" || statusStr == "online" {
		status = models.StatusOnline
		health = models.HealthHealthy
	} else if statusStr == "stopped" || statusStr == "offline" {
		status = models.StatusOffline
		health = models.HealthUnknown
	} else {
		status = models.StatusUnknown
		health = models.HealthUnknown
	}

	serverStatus := models.NewServerStatus(
		health,
		status,
		"Minecraft server status: "+statusStr,
		"0h 30m", // Would be calculated from actual uptime
		20.0,     // Would be actual TPS from server
	)

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
	output, err := (*ctx.MinecraftServer).ExecuteCommand(command.Command)
	if err != nil {
		slog.Error("Failed to execute command", "command", command.Command, "error", err.Error())
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to execute command: "+err.Error()))
		return
	}

	// Return command result
	result := models.NewCommandResultWithOutput(output)

	slog.Info("Successfully executed command", "command", command.Command)
	c.JSON(http.StatusOK, result)
}
