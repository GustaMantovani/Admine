package handlers

import (
	"net/http"

	"github.com/GustaMantovani/Admine/server_handler/internal"
	"github.com/GustaMantovani/Admine/server_handler/internal/api/models"
	"github.com/GustaMantovani/Admine/server_handler/pkg"
	"github.com/gin-gonic/gin"
)

// ServerHandler handles server-related API endpoints
type ServerHandler struct {
}

// NewServerHandler creates a new ServerHandler instance
func NewServerHandler() *ServerHandler {
	return &ServerHandler{}
}

// GetInfo handles GET /info endpoint
func (h *ServerHandler) GetInfo(c *gin.Context) {
	pkg.Logger.Info("GET /info endpoint called")

	ctx := internal.Get()
	if ctx.MinecraftServer == nil {
		pkg.Logger.Error("MinecraftServer is not initialized")
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Minecraft server not initialized"))
		return
	}

	// Get server info through the MinecraftServer interface
	_, err := (*ctx.MinecraftServer).Info()
	if err != nil {
		pkg.Logger.Error("Failed to get server info: %s", err.Error())
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

	pkg.Logger.Info("Successfully retrieved server info")
	c.JSON(http.StatusOK, serverInfo)
}

// GetStatus handles GET /status endpoint
func (h *ServerHandler) GetStatus(c *gin.Context) {
	pkg.Logger.Info("GET /status endpoint called")

	ctx := internal.Get()
	if ctx.MinecraftServer == nil {
		pkg.Logger.Error("MinecraftServer is not initialized")
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Minecraft server not initialized"))
		return
	}

	// Get server status through the MinecraftServer interface
	statusStr, err := (*ctx.MinecraftServer).Status()
	if err != nil {
		pkg.Logger.Error("Failed to get server status: %s", err.Error())
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

	pkg.Logger.Info("Successfully retrieved server status")
	c.JSON(http.StatusOK, serverStatus)
}

// PostCommand handles POST /command endpoint
func (h *ServerHandler) PostCommand(c *gin.Context) {
	pkg.Logger.Info("POST /command endpoint called")

	var command models.Command
	if err := c.ShouldBindJSON(&command); err != nil {
		pkg.Logger.Error("Invalid command request: %s", err.Error())
		c.JSON(http.StatusBadRequest, models.NewErrorResponse("Invalid request: "+err.Error()))
		return
	}

	ctx := internal.Get()
	if ctx.MinecraftServer == nil {
		pkg.Logger.Error("MinecraftServer is not initialized")
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Minecraft server not initialized"))
		return
	}

	// Execute command through the MinecraftServer interface
	output, err := (*ctx.MinecraftServer).ExecuteCommand(command.Command)
	if err != nil {
		pkg.Logger.Error("Failed to execute command '%s': %s", command.Command, err.Error())
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to execute command: "+err.Error()))
		return
	}

	// Return command result
	result := models.NewCommandResultWithOutput(output)

	pkg.Logger.Info("Successfully executed command '%s'", command.Command)
	c.JSON(http.StatusOK, result)
}
