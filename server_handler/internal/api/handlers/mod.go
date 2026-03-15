package handlers

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	apimodels "github.com/GustaMantovani/Admine/server_handler/internal/api/models"
	"github.com/GustaMantovani/Admine/server_handler/internal/config"
	"github.com/GustaMantovani/Admine/server_handler/internal/pubsub"
	"github.com/GustaMantovani/Admine/server_handler/internal/server"
	"github.com/gin-gonic/gin"
)

// ModHandler handles mod-related API endpoints
type ModHandler struct {
	server        server.MinecraftServer
	pubsub        pubsub.PubSubService
	origin        string
	serverChannel string
	modTimeout    time.Duration
	mainCtx       context.Context
}

// NewModHandler creates a new ModHandler
func NewModHandler(
	srv server.MinecraftServer,
	ps pubsub.PubSubService,
	origin string,
	serverChannel string,
	cfg config.MinecraftServerConfig,
	mainCtx context.Context,
) *ModHandler {
	modTimeout := cfg.ModInstallTimeout
	if modTimeout == 0 {
		modTimeout = 2 * time.Minute
	}

	return &ModHandler{
		server:        srv,
		pubsub:        ps,
		origin:        origin,
		serverChannel: serverChannel,
		modTimeout:    modTimeout,
		mainCtx:       mainCtx,
	}
}

func (h *ModHandler) publish(tags []string, message string) {
	msg := pubsub.NewAdmineMessage(h.origin, tags, message)
	h.pubsub.Publish(h.serverChannel, msg)
}

// PostInstallMod handles POST /mods
// Accepts multipart/form-data (file upload) or application/json (URL download)
func (h *ModHandler) PostInstallMod(c *gin.Context) {
	slog.Info("POST /mods endpoint called")

	if h.server == nil {
		slog.Error("MinecraftServer is not initialized")
		c.JSON(http.StatusInternalServerError, apimodels.NewModInstallResponse("error", "Minecraft server not initialized"))
		return
	}

	if strings.HasPrefix(c.ContentType(), "multipart/form-data") {
		h.handleFileUpload(c)
	} else {
		h.handleURLDownload(c)
	}
}

func (h *ModHandler) handleFileUpload(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		slog.Error("Failed to get uploaded file", "error", err)
		c.JSON(http.StatusBadRequest, apimodels.NewModInstallResponse("error", "Failed to get uploaded file: "+err.Error()))
		return
	}
	defer file.Close()

	fileName := header.Filename
	if !isJarFile(fileName) {
		c.JSON(http.StatusBadRequest, apimodels.NewModInstallResponse("error", "Invalid file type: only .jar files are accepted"))
		return
	}

	fileData, err := io.ReadAll(file)
	if err != nil {
		slog.Error("Failed to read uploaded file", "error", err)
		c.JSON(http.StatusInternalServerError, apimodels.NewModInstallResponse("error", "Failed to read file data"))
		return
	}

	c.JSON(http.StatusAccepted, apimodels.NewModInstallResponse("accepted", "Mod installation started for: "+fileName))

	go h.installMod(fileName, fileData)
}

func (h *ModHandler) handleURLDownload(c *gin.Context) {
	var req apimodels.ModInstallRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Error("Invalid mod install request", "error", err)
		c.JSON(http.StatusBadRequest, apimodels.NewModInstallResponse("error", "Invalid request: "+err.Error()))
		return
	}

	fileName := filepath.Base(req.URL)
	if !isJarFile(fileName) {
		c.JSON(http.StatusBadRequest, apimodels.NewModInstallResponse("error", "Invalid file type: URL must point to a .jar file"))
		return
	}

	c.JSON(http.StatusAccepted, apimodels.NewModInstallResponse("accepted", "Mod download and installation started for: "+fileName))

	go h.downloadAndInstallMod(req.URL, fileName)
}

func (h *ModHandler) installMod(fileName string, fileData []byte) {
	ctx, cancel := context.WithTimeout(h.mainCtx, h.modTimeout)
	defer cancel()

	h.publish([]string{"notification"}, "Installing mod: "+fileName)

	reader := strings.NewReader(string(fileData))
	result, err := h.server.InstallMod(ctx, fileName, reader)
	if err != nil {
		slog.Error("Failed to install mod", "file", fileName, "error", err)
		if ctx.Err() != nil {
			return
		}
		h.publish([]string{"mod_install_result"}, "Failed to install mod "+fileName+": "+err.Error())
		return
	}

	h.publish([]string{"mod_install_result"}, result.Message)
	slog.Info("Mod installed successfully", "file", fileName)
}

func (h *ModHandler) downloadAndInstallMod(url string, fileName string) {
	ctx, cancel := context.WithTimeout(h.mainCtx, h.modTimeout)
	defer cancel()

	h.publish([]string{"notification"}, "Downloading mod: "+fileName)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		slog.Error("Failed to create download request", "url", url, "error", err)
		if ctx.Err() != nil {
			return
		}
		h.publish([]string{"mod_install_result"}, "Failed to download mod "+fileName+": "+err.Error())
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		slog.Error("Failed to download mod", "url", url, "error", err)
		if ctx.Err() != nil {
			return
		}
		h.publish([]string{"mod_install_result"}, "Failed to download mod "+fileName+": "+err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.Error("Failed to download mod", "url", url, "status", resp.StatusCode)
		h.publish([]string{"mod_install_result"}, "Failed to download mod "+fileName+": HTTP "+resp.Status)
		return
	}

	h.publish([]string{"notification"}, "Download complete, installing mod: "+fileName)

	result, err := h.server.InstallMod(ctx, fileName, resp.Body)
	if err != nil {
		slog.Error("Failed to install mod", "file", fileName, "error", err)
		if ctx.Err() != nil {
			return
		}
		h.publish([]string{"mod_install_result"}, "Failed to install mod "+fileName+": "+err.Error())
		return
	}

	h.publish([]string{"mod_install_result"}, result.Message)
	slog.Info("Mod downloaded and installed successfully", "file", fileName, "url", url)
}

// GetListMods handles GET /mods
func (h *ModHandler) GetListMods(c *gin.Context) {
	slog.Info("GET /mods endpoint called")

	if h.server == nil {
		slog.Error("MinecraftServer is not initialized")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Minecraft server not initialized"})
		return
	}

	result, err := h.server.ListMods(c.Request.Context())
	if err != nil {
		slog.Error("Failed to list mods", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list mods: " + err.Error()})
		return
	}

	slog.Info("Mods listed successfully", "count", result.Total)
	c.JSON(http.StatusOK, result)
}

// DeleteRemoveMod handles DELETE /mods/:filename
func (h *ModHandler) DeleteRemoveMod(c *gin.Context) {
	fileName := c.Param("filename")
	slog.Info("DELETE /mods endpoint called", "filename", fileName)

	if !isJarFile(fileName) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type: only .jar files can be removed"})
		return
	}

	if h.server == nil {
		slog.Error("MinecraftServer is not initialized")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Minecraft server not initialized"})
		return
	}

	result, err := h.server.RemoveMod(c.Request.Context(), fileName)
	if err != nil {
		slog.Error("Failed to remove mod", "filename", fileName, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove mod: " + err.Error()})
		return
	}

	slog.Info("Mod removed successfully", "filename", fileName)
	c.JSON(http.StatusOK, result)
}

func isJarFile(fileName string) bool {
	return strings.HasSuffix(strings.ToLower(fileName), ".jar")
}
