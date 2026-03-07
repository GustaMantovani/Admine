package handlers

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/GustaMantovani/Admine/server_handler/internal"
	apimodels "github.com/GustaMantovani/Admine/server_handler/internal/api/models"
	"github.com/GustaMantovani/Admine/server_handler/internal/pubsub"
	pubsubmodels "github.com/GustaMantovani/Admine/server_handler/internal/pubsub/models"
	"github.com/gin-gonic/gin"
)

// ModHandler handles mod-related API endpoints
type ModHandler struct {
	pubsub pubsub.PubSubService
}

// NewModHandler creates a new ModHandler instance
func NewModHandler(ps pubsub.PubSubService) *ModHandler {
	return &ModHandler{
		pubsub: ps,
	}
}

// PostInstallMod handles POST /mods endpoint
// Accepts either multipart/form-data (file upload) or application/json (URL download)
func (h *ModHandler) PostInstallMod(c *gin.Context) {
	slog.Info("POST /mods endpoint called")

	appContext := internal.Get()
	if appContext.MinecraftServer == nil {
		slog.Error("MinecraftServer is not initialized")
		c.JSON(http.StatusInternalServerError, apimodels.NewModInstallResponse("error", "Minecraft server not initialized"))
		return
	}

	contentType := c.ContentType()

	if strings.HasPrefix(contentType, "multipart/form-data") {
		h.handleFileUpload(c, appContext)
	} else {
		h.handleURLDownload(c, appContext)
	}
}

// handleFileUpload processes a multipart file upload
func (h *ModHandler) handleFileUpload(c *gin.Context, appContext *internal.AppContext) {
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

	// Read file data into memory for async processing
	fileData, err := io.ReadAll(file)
	if err != nil {
		slog.Error("Failed to read uploaded file", "error", err)
		c.JSON(http.StatusInternalServerError, apimodels.NewModInstallResponse("error", "Failed to read file data"))
		return
	}

	// Return 202 Accepted immediately
	c.JSON(http.StatusAccepted, apimodels.NewModInstallResponse("accepted", "Mod installation started for: "+fileName))

	// Process asynchronously
	go h.installMod(appContext, fileName, fileData)
}

// handleURLDownload processes a JSON request with a URL
func (h *ModHandler) handleURLDownload(c *gin.Context, appContext *internal.AppContext) {
	var req apimodels.ModInstallRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Error("Invalid mod install request", "error", err)
		c.JSON(http.StatusBadRequest, apimodels.NewModInstallResponse("error", "Invalid request: "+err.Error()))
		return
	}

	// Extract filename from URL
	fileName := filepath.Base(req.URL)
	if !isJarFile(fileName) {
		c.JSON(http.StatusBadRequest, apimodels.NewModInstallResponse("error", "Invalid file type: URL must point to a .jar file"))
		return
	}

	// Return 202 Accepted immediately
	c.JSON(http.StatusAccepted, apimodels.NewModInstallResponse("accepted", "Mod download and installation started for: "+fileName))

	// Process asynchronously: download then install
	go h.downloadAndInstallMod(appContext, req.URL, fileName)
}

// installMod installs a mod from in-memory data
func (h *ModHandler) installMod(appContext *internal.AppContext, fileName string, fileData []byte) {
	modInstallTimeout := appContext.Config.MinecraftServer.ModInstallTimeout
	if modInstallTimeout == 0 {
		modInstallTimeout = 2 * time.Minute
	}

	ctx, cancel := context.WithTimeout(*appContext.MainCtx, modInstallTimeout)
	defer cancel()

	// Notify start
	startMsg := pubsubmodels.NewAdmineMessage([]string{"notification"}, "Installing mod: "+fileName)
	h.pubsub.Publish(appContext.Config.PubSub.AdmineChannelsMap.ServerChannel, startMsg)

	reader := strings.NewReader(string(fileData))
	result, err := (*appContext.MinecraftServer).InstallMod(ctx, fileName, reader)
	if err != nil {
		slog.Error("Failed to install mod", "file", fileName, "error", err)
		if ctx.Err() != nil {
			return
		}
		errorMsg := pubsubmodels.NewAdmineMessage([]string{"mod_install_result"}, "Failed to install mod "+fileName+": "+err.Error())
		h.pubsub.Publish(appContext.Config.PubSub.AdmineChannelsMap.ServerChannel, errorMsg)
		return
	}

	successMsg := pubsubmodels.NewAdmineMessage([]string{"mod_install_result"}, result.Message)
	h.pubsub.Publish(appContext.Config.PubSub.AdmineChannelsMap.ServerChannel, successMsg)
	slog.Info("Mod installed successfully", "file", fileName)
}

// downloadAndInstallMod downloads a mod from URL and installs it
func (h *ModHandler) downloadAndInstallMod(appContext *internal.AppContext, url string, fileName string) {
	modInstallTimeout := appContext.Config.MinecraftServer.ModInstallTimeout
	if modInstallTimeout == 0 {
		modInstallTimeout = 2 * time.Minute
	}

	ctx, cancel := context.WithTimeout(*appContext.MainCtx, modInstallTimeout)
	defer cancel()

	// Notify start
	startMsg := pubsubmodels.NewAdmineMessage([]string{"notification"}, "Downloading mod: "+fileName)
	h.pubsub.Publish(appContext.Config.PubSub.AdmineChannelsMap.ServerChannel, startMsg)

	// Download the mod
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		slog.Error("Failed to create download request", "url", url, "error", err)
		if ctx.Err() != nil {
			return
		}
		errorMsg := pubsubmodels.NewAdmineMessage([]string{"mod_install_result"}, "Failed to download mod "+fileName+": "+err.Error())
		h.pubsub.Publish(appContext.Config.PubSub.AdmineChannelsMap.ServerChannel, errorMsg)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		slog.Error("Failed to download mod", "url", url, "error", err)
		if ctx.Err() != nil {
			return
		}
		errorMsg := pubsubmodels.NewAdmineMessage([]string{"mod_install_result"}, "Failed to download mod "+fileName+": "+err.Error())
		h.pubsub.Publish(appContext.Config.PubSub.AdmineChannelsMap.ServerChannel, errorMsg)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.Error("Failed to download mod", "url", url, "status", resp.StatusCode)
		errorMsg := pubsubmodels.NewAdmineMessage([]string{"mod_install_result"}, "Failed to download mod "+fileName+": HTTP "+resp.Status)
		h.pubsub.Publish(appContext.Config.PubSub.AdmineChannelsMap.ServerChannel, errorMsg)
		return
	}

	// Notify download complete, starting install
	downloadMsg := pubsubmodels.NewAdmineMessage([]string{"notification"}, "Download complete, installing mod: "+fileName)
	h.pubsub.Publish(appContext.Config.PubSub.AdmineChannelsMap.ServerChannel, downloadMsg)

	result, err := (*appContext.MinecraftServer).InstallMod(ctx, fileName, resp.Body)
	if err != nil {
		slog.Error("Failed to install mod", "file", fileName, "error", err)
		if ctx.Err() != nil {
			return
		}
		errorMsg := pubsubmodels.NewAdmineMessage([]string{"mod_install_result"}, "Failed to install mod "+fileName+": "+err.Error())
		h.pubsub.Publish(appContext.Config.PubSub.AdmineChannelsMap.ServerChannel, errorMsg)
		return
	}

	successMsg := pubsubmodels.NewAdmineMessage([]string{"mod_install_result"}, result.Message)
	h.pubsub.Publish(appContext.Config.PubSub.AdmineChannelsMap.ServerChannel, successMsg)
	slog.Info("Mod downloaded and installed successfully", "file", fileName, "url", url)
}

// isJarFile checks if the filename has a .jar extension
func isJarFile(fileName string) bool {
	return strings.HasSuffix(strings.ToLower(fileName), ".jar")
}

// GetListMods handles GET /mods endpoint
func (h *ModHandler) GetListMods(c *gin.Context) {
	slog.Info("GET /mods endpoint called")

	appContext := internal.Get()
	if appContext.MinecraftServer == nil {
		slog.Error("MinecraftServer is not initialized")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Minecraft server not initialized"})
		return
	}

	result, err := (*appContext.MinecraftServer).ListMods(c.Request.Context())
	if err != nil {
		slog.Error("Failed to list mods", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list mods: " + err.Error()})
		return
	}

	slog.Info("Mods listed successfully", "count", result.Total)
	c.JSON(http.StatusOK, result)
}

// DeleteRemoveMod handles DELETE /mods/:filename endpoint
func (h *ModHandler) DeleteRemoveMod(c *gin.Context) {
	fileName := c.Param("filename")
	slog.Info("DELETE /mods endpoint called", "filename", fileName)

	if !isJarFile(fileName) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type: only .jar files can be removed"})
		return
	}

	appContext := internal.Get()
	if appContext.MinecraftServer == nil {
		slog.Error("MinecraftServer is not initialized")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Minecraft server not initialized"})
		return
	}

	result, err := (*appContext.MinecraftServer).RemoveMod(c.Request.Context(), fileName)
	if err != nil {
		slog.Error("Failed to remove mod", "filename", fileName, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove mod: " + err.Error()})
		return
	}

	slog.Info("Mod removed successfully", "filename", fileName)
	c.JSON(http.StatusOK, result)
}
