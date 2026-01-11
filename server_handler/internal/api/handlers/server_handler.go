package handlers

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/GustaMantovani/Admine/server_handler/internal"
	"github.com/GustaMantovani/Admine/server_handler/internal/api/models"
	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
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

// GetResourceUsage handles GET /resources endpoint
func (h *ServerHandler) GetResourceUsage(c *gin.Context) {
	slog.Info("GET /resources endpoint called")

	cpuPercentages, err := cpu.Percent(time.Second, false)
	if err != nil {
		slog.Error("Failed to get CPU usage", "error", err.Error())
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to get CPU usage: "+err.Error()))
		return
	}

	memInfo, err := mem.VirtualMemory()
	if err != nil {
		slog.Error("Failed to get memory usage", "error", err.Error())
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to get memory usage: "+err.Error()))
		return
	}

	diskInfo, err := disk.Usage("/")
	if err != nil {
		slog.Error("Failed to get disk usage", "error", err.Error())
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to get disk usage: "+err.Error()))
		return
	}

	resourceUsage := models.ResourceUsage{
		CPUUsage:        cpuPercentages[0],
		MemoryUsed:      memInfo.Used,
		MemoryTotal:     memInfo.Total,
		MemoryUsedPercent: memInfo.UsedPercent,
		DiskUsed:        diskInfo.Used,
		DiskTotal:       diskInfo.Total,
		DiskUsedPercent: diskInfo.UsedPercent,
	}

	slog.Info("Successfully retrieved resource usage")
	c.JSON(http.StatusOK, resourceUsage)
}
