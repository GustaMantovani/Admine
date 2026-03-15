package handlers

import (
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/GustaMantovani/Admine/server_handler/internal/api/models"
	"github.com/GustaMantovani/Admine/server_handler/internal/server"
	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
)

const (
	defaultLogLines = 100
	maxLogLines     = 100
)

// ServerHandler handles server-related API endpoints
type ServerHandler struct {
	server server.MinecraftServer
}

// NewServerHandler creates a new ServerHandler
func NewServerHandler(srv server.MinecraftServer) *ServerHandler {
	return &ServerHandler{server: srv}
}

// GetInfo handles GET /info
func (h *ServerHandler) GetInfo(c *gin.Context) {
	slog.Info("GET /info endpoint called")

	if h.server == nil {
		slog.Error("MinecraftServer is not initialized")
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Minecraft server not initialized"))
		return
	}

	serverInfo, err := h.server.Info(c.Request.Context())
	if err != nil {
		slog.Error("Failed to get server info", "error", err.Error())
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to get server info: "+err.Error()))
		return
	}

	slog.Info("Successfully retrieved server info")
	c.JSON(http.StatusOK, serverInfo)
}

// GetStatus handles GET /status
func (h *ServerHandler) GetStatus(c *gin.Context) {
	slog.Info("GET /status endpoint called")

	if h.server == nil {
		slog.Error("MinecraftServer is not initialized")
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Minecraft server not initialized"))
		return
	}

	serverStatus, err := h.server.Status(c.Request.Context())
	if err != nil {
		slog.Error("Failed to get server status", "error", err.Error())
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to get server status: "+err.Error()))
		return
	}

	slog.Info("Successfully retrieved server status")
	c.JSON(http.StatusOK, serverStatus)
}

// GetLogs handles GET /logs?n=<int>
func (h *ServerHandler) GetLogs(c *gin.Context) {
	slog.Info("GET /logs endpoint called")

	if h.server == nil {
		slog.Error("MinecraftServer is not initialized")
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Minecraft server not initialized"))
		return
	}

	n := defaultLogLines
	if nRaw := c.Query("n"); nRaw != "" {
		parsedN, err := strconv.Atoi(nRaw)
		if err != nil {
			slog.Error("Invalid query param n", "n", nRaw, "error", err.Error())
			c.JSON(http.StatusBadRequest, models.NewErrorResponse("Invalid query param 'n': must be an integer between 1 and 100"))
			return
		}
		n = parsedN
	}

	if n < 1 || n > maxLogLines {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse("Invalid query param 'n': must be an integer between 1 and 100"))
		return
	}

	logs, err := h.server.Logs(c.Request.Context(), n)
	if err != nil {
		slog.Error("Failed to get server logs", "error", err.Error())
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to get server logs: "+err.Error()))
		return
	}

	slog.Info("Successfully retrieved server logs", "lines", len(logs))
	c.JSON(http.StatusOK, models.NewLogsResponse(logs))
}

// PostCommand handles POST /command
func (h *ServerHandler) PostCommand(c *gin.Context) {
	slog.Info("POST /command endpoint called")

	var command models.Command
	if err := c.ShouldBindJSON(&command); err != nil {
		slog.Error("Invalid command request", "error", err.Error())
		c.JSON(http.StatusBadRequest, models.NewErrorResponse("Invalid request: "+err.Error()))
		return
	}

	if h.server == nil {
		slog.Error("MinecraftServer is not initialized")
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Minecraft server not initialized"))
		return
	}

	result, err := h.server.ExecuteCommand(c.Request.Context(), command.Command)
	if err != nil {
		slog.Error("Failed to execute command", "command", command.Command, "error", err.Error())
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to execute command: "+err.Error()))
		return
	}

	slog.Info("Successfully executed command", "command", command.Command)
	c.JSON(http.StatusOK, result)
}

// GetResourceUsage handles GET /resources
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
		CPUUsage:          cpuPercentages[0],
		MemoryUsed:        memInfo.Used,
		MemoryTotal:       memInfo.Total,
		MemoryUsedPercent: memInfo.UsedPercent,
		DiskUsed:          diskInfo.Used,
		DiskTotal:         diskInfo.Total,
		DiskUsedPercent:   diskInfo.UsedPercent,
	}

	slog.Info("Successfully retrieved resource usage")
	c.JSON(http.StatusOK, resourceUsage)
}
