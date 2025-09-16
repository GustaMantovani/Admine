package pubsub

import (
	"log/slog"
	"strings"

	"github.com/GustaMantovani/Admine/server_handler/internal"
	"github.com/GustaMantovani/Admine/server_handler/internal/pubsub/models"
)

// EventHandler handles incoming messages and manages server operations
type EventHandler struct {
	pubsub PubSubService
}

// NewEventHandler creates a new EventHandler instance
func NewEventHandler(ps PubSubService) *EventHandler {
	return &EventHandler{
		pubsub: ps,
	}
}

// ManageCommand processes incoming messages and routes them to appropriate handlers
func (eh *EventHandler) ManageCommand(msg *models.AdmineMessage) error {
	if msg.HasTag("server_on") {
		eh.serverUp()
	} else if msg.HasTag("server_off") {
		eh.serverOff()
	} else if msg.HasTag("restart") {
		eh.restart()
	} else if msg.HasTag("command") {
		eh.command(msg.Message)
	} else {
		ctx := internal.Get()
		responseMsg := models.NewAdmineMessage([]string{"error"}, "Invalid tag.")
		eh.pubsub.Publish(ctx.Config.PubSub.AdmineChannelsMap.ServerChannel, responseMsg)

		slog.Error("Received an invalid tag", "tags", msg.Tags)
	}

	return nil
}

func (eh *EventHandler) serverUp() {
	ctx := internal.Get()
	// Start the server using the MinecraftServer interface
	if ctx.MinecraftServer == nil {
		slog.Error("MinecraftServer is not initialized")
		responseMsg := models.NewAdmineMessage([]string{"error"}, "Server not initialized")
		eh.pubsub.Publish(ctx.Config.PubSub.AdmineChannelsMap.ServerChannel, responseMsg)
		return
	}

	// Send starting message
	responseMsg := models.NewAdmineMessage([]string{"server_status"}, "Starting server")
	eh.pubsub.Publish(ctx.Config.PubSub.AdmineChannelsMap.ServerChannel, responseMsg)

	err := (*ctx.MinecraftServer).Start()
	if err != nil {
		slog.Error("Error starting server", "error", err.Error())
		errorMsg := models.NewAdmineMessage([]string{"error"}, "Failed to start server: "+err.Error())
		eh.pubsub.Publish(ctx.Config.PubSub.AdmineChannelsMap.ServerChannel, errorMsg)
		return
	}

	// Get server startInfo (this might include network information)
	startInfo := (*ctx.MinecraftServer).StartUpInfo()

	successMsg := models.NewAdmineMessage([]string{"server_on"}, startInfo)
	eh.pubsub.Publish(ctx.Config.PubSub.AdmineChannelsMap.ServerChannel, successMsg)

	slog.Info("Server started successfully")
}

func (eh *EventHandler) serverOff() {
	ctx := internal.Get()
	if ctx.MinecraftServer == nil {
		slog.Error("MinecraftServer is not initialized")
		responseMsg := models.NewAdmineMessage([]string{"error"}, "Server not initialized")
		eh.pubsub.Publish(ctx.Config.PubSub.AdmineChannelsMap.ServerChannel, responseMsg)
		return
	}

	// Send stopping message
	responseMsg := models.NewAdmineMessage([]string{"server_status"}, "Stopping server")
	eh.pubsub.Publish(ctx.Config.PubSub.AdmineChannelsMap.ServerChannel, responseMsg)

	// Execute stop command through the server interface
	_, err := (*ctx.MinecraftServer).ExecuteCommand("/stop")
	if err != nil {
		slog.Error("Error executing stop command", "error", err.Error())
	}

	// Wait for graceful shutdown (simplified approach)
	// In a real implementation, you might want to monitor server logs
	// or implement a more sophisticated shutdown detection

	err = (*ctx.MinecraftServer).Stop()
	if err != nil {
		slog.Error("Error stopping server", "error", err.Error())
		errorMsg := models.NewAdmineMessage([]string{"error"}, "Error stopping server: "+err.Error())
		eh.pubsub.Publish(ctx.Config.PubSub.AdmineChannelsMap.ServerChannel, errorMsg)
		return
	}

	successMsg := models.NewAdmineMessage([]string{"server_off"}, "Server stopped successfully")
	eh.pubsub.Publish(ctx.Config.PubSub.AdmineChannelsMap.ServerChannel, successMsg)

	slog.Info("Server stopped successfully")
}

func (eh *EventHandler) restart() {
	ctx := internal.Get()
	if ctx.MinecraftServer == nil {
		slog.Error("MinecraftServer is not initialized")
		responseMsg := models.NewAdmineMessage([]string{"error"}, "Server not initialized")
		eh.pubsub.Publish(ctx.Config.PubSub.AdmineChannelsMap.ServerChannel, responseMsg)
		return
	}

	// Send restarting message
	statusMsg := models.NewAdmineMessage([]string{"server_status"}, "Restarting server")
	eh.pubsub.Publish(ctx.Config.PubSub.AdmineChannelsMap.ServerChannel, statusMsg)

	slog.Info("Starting server restart process")

	// Use the Restart method from the MinecraftServer interface
	err := (*ctx.MinecraftServer).Restart()
	if err != nil {
		slog.Error("Error restarting server", "error", err.Error())
		errorMsg := models.NewAdmineMessage([]string{"error"}, "Failed to restart server: "+err.Error())
		eh.pubsub.Publish(ctx.Config.PubSub.AdmineChannelsMap.ServerChannel, errorMsg)
		return
	}

	// Get server startup info after restart
	startInfo := (*ctx.MinecraftServer).StartUpInfo()

	// Send restart success message
	successMsg := models.NewAdmineMessage([]string{"restart_complete"}, startInfo)
	eh.pubsub.Publish(ctx.Config.PubSub.AdmineChannelsMap.ServerChannel, successMsg)

	slog.Info("Server restarted successfully")
}

func (eh *EventHandler) command(message string) {
	ctx := internal.Get()
	if ctx.MinecraftServer == nil {
		slog.Error("MinecraftServer is not initialized")
		responseMsg := models.NewAdmineMessage([]string{"error"}, "Server not initialized")
		eh.pubsub.Publish(ctx.Config.PubSub.AdmineChannelsMap.ServerChannel, responseMsg)
		return
	}

	// Execute the command through the MinecraftServer interface
	result, err := (*ctx.MinecraftServer).ExecuteCommand(message)
	if err != nil {
		slog.Error("Error executing command", "command", message, "error", err.Error())
		errorMsg := models.NewAdmineMessage([]string{"error"}, "Failed to execute command: "+err.Error())
		eh.pubsub.Publish(ctx.Config.PubSub.AdmineChannelsMap.ServerChannel, errorMsg)
		return
	}

	// Send command result
	var responseMessage string
	if strings.TrimSpace(result) != "" {
		responseMessage = result
	} else {
		responseMessage = "Command executed successfully"
	}

	successMsg := models.NewAdmineMessage([]string{"command_result"}, responseMessage)
	eh.pubsub.Publish(ctx.Config.PubSub.AdmineChannelsMap.ServerChannel, successMsg)

	slog.Info("Executed command successfully", "command", message)
}
