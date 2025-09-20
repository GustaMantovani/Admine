package pubsub

import (
	"context"
	"log/slog"
	"strings"
	"time"

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
	} else if msg.HasTag("server_down") {
		eh.serverDown()
	} else if msg.HasTag("restart") {
		eh.restart()
	} else if msg.HasTag("command") {
		eh.command(msg.Message)
	} else {
		ctx := internal.Get()
		responseMsg := models.NewAdmineMessage([]string{"notification"}, "Invalid tag.")
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
		responseMsg := models.NewAdmineMessage([]string{"notification"}, "Server not initialized")
		eh.pubsub.Publish(ctx.Config.PubSub.AdmineChannelsMap.ServerChannel, responseMsg)
		return
	}

	// Send starting message
	responseMsg := models.NewAdmineMessage([]string{"notification"}, "Starting server")
	eh.pubsub.Publish(ctx.Config.PubSub.AdmineChannelsMap.ServerChannel, responseMsg)

	// Create context with timeout for server start operation
	startCtx, cancel := context.WithTimeout(*ctx.MainCtx, ctx.Config.MinecraftServer.ServerOnTimeout)
	defer cancel()

	err := (*ctx.MinecraftServer).Start(startCtx)
	if err != nil {
		slog.Error("Error starting server", "error", err.Error())
		errorMsg := models.NewAdmineMessage([]string{"notification"}, "Failed to start server: "+err.Error())
		eh.pubsub.Publish(ctx.Config.PubSub.AdmineChannelsMap.ServerChannel, errorMsg)
		return
	}

	// Get server startInfo (this might include network information)
	infoCtx, infoCancel := context.WithTimeout(*ctx.MainCtx, 30*time.Second)
	defer infoCancel()

	startInfo := (*ctx.MinecraftServer).StartUpInfo(infoCtx)

	successMsg := models.NewAdmineMessage([]string{"server_on"}, startInfo)
	eh.pubsub.Publish(ctx.Config.PubSub.AdmineChannelsMap.ServerChannel, successMsg)

	slog.Info("Server started successfully")
}

func (eh *EventHandler) serverOff() {
	ctx := internal.Get()
	if ctx.MinecraftServer == nil {
		slog.Error("MinecraftServer is not initialized")
		responseMsg := models.NewAdmineMessage([]string{"notification"}, "Server not initialized")
		eh.pubsub.Publish(ctx.Config.PubSub.AdmineChannelsMap.ServerChannel, responseMsg)
		return
	}

	// Send stopping message
	responseMsg := models.NewAdmineMessage([]string{"notification"}, "Stopping server")
	eh.pubsub.Publish(ctx.Config.PubSub.AdmineChannelsMap.ServerChannel, responseMsg)

	// Create context with timeout for server stop operation
	stopCtx, cancel := context.WithTimeout(*ctx.MainCtx, ctx.Config.MinecraftServer.ServerOffTimeout)
	defer cancel()

	err := (*ctx.MinecraftServer).Stop(stopCtx)
	if err != nil {
		slog.Error("Error stopping server", "error", err.Error())
		errorMsg := models.NewAdmineMessage([]string{"notification"}, "Error stopping server: "+err.Error())
		eh.pubsub.Publish(ctx.Config.PubSub.AdmineChannelsMap.ServerChannel, errorMsg)
		return
	}

	successMsg := models.NewAdmineMessage([]string{"server_off"}, "Server stopped successfully")
	eh.pubsub.Publish(ctx.Config.PubSub.AdmineChannelsMap.ServerChannel, successMsg)

	slog.Info("Server stopped successfully")
}

func (eh *EventHandler) serverDown() {
	ctx := internal.Get()
	if ctx.MinecraftServer == nil {
		slog.Error("MinecraftServer is not initialized")
		responseMsg := models.NewAdmineMessage([]string{"notification"}, "Server not initialized")
		eh.pubsub.Publish(ctx.Config.PubSub.AdmineChannelsMap.ServerChannel, responseMsg)
		return
	}

	// Send stopping message
	responseMsg := models.NewAdmineMessage([]string{"notification"}, "Removing server")
	eh.pubsub.Publish(ctx.Config.PubSub.AdmineChannelsMap.ServerChannel, responseMsg)

	// Execute stop command through the server interface with timeout
	cmdCtx, cmdCancel := context.WithTimeout(*ctx.MainCtx, ctx.Config.MinecraftServer.ServerCommandExecTimeout)
	defer cmdCancel()

	_, err := (*ctx.MinecraftServer).ExecuteCommand(cmdCtx, "/stop")
	if err != nil {
		slog.Error("Error executing stop command", "error", err.Error())
	}

	// Create context with timeout for server down operation
	downCtx, downCancel := context.WithTimeout(*ctx.MainCtx, ctx.Config.MinecraftServer.ServerOffTimeout)
	defer downCancel()

	err = (*ctx.MinecraftServer).Down(downCtx)
	if err != nil {
		slog.Error("Error stopping server", "error", err.Error())
		errorMsg := models.NewAdmineMessage([]string{"notification"}, "Error removing server: "+err.Error())
		eh.pubsub.Publish(ctx.Config.PubSub.AdmineChannelsMap.ServerChannel, errorMsg)
		return
	}

	successMsg := models.NewAdmineMessage([]string{"server_off"}, "Server removed successfully")
	eh.pubsub.Publish(ctx.Config.PubSub.AdmineChannelsMap.ServerChannel, successMsg)

	slog.Info("Server removed successfully")
}

func (eh *EventHandler) restart() {
	ctx := internal.Get()

	if ctx.MinecraftServer == nil {
		slog.Error("MinecraftServer is not initialized")
		responseMsg := models.NewAdmineMessage([]string{"notification"}, "Server not initialized")
		eh.pubsub.Publish(ctx.Config.PubSub.AdmineChannelsMap.ServerChannel, responseMsg)
		return
	}

	// Send restarting message
	statusMsg := models.NewAdmineMessage([]string{"notification"}, "Restarting server")
	eh.pubsub.Publish(ctx.Config.PubSub.AdmineChannelsMap.ServerChannel, statusMsg)

	slog.Info("Starting server restart process")

	// Create context with timeout for command execution
	cmdCtx, cancel := context.WithTimeout(*ctx.MainCtx, ctx.Config.MinecraftServer.ServerCommandExecTimeout)
	defer cancel()

	// Execute stop command and let Docker restart the container
	_, err := (*ctx.MinecraftServer).ExecuteCommand(cmdCtx, "/stop")
	if err != nil {
		slog.Error("Error executing stop command for restart", "error", err.Error())
		errorMsg := models.NewAdmineMessage([]string{"notification"}, "Failed to restart server: "+err.Error())
		eh.pubsub.Publish(ctx.Config.PubSub.AdmineChannelsMap.ServerChannel, errorMsg)
		return
	}

	slog.Info("Server restarted successfully")
}

func (eh *EventHandler) command(message string) {
	appContext := internal.Get()
	if appContext.MinecraftServer == nil {
		slog.Error("MinecraftServer is not initialized")
		responseMsg := models.NewAdmineMessage([]string{"notification"}, "Server not initialized")
		eh.pubsub.Publish(appContext.Config.PubSub.AdmineChannelsMap.ServerChannel, responseMsg)
		return
	}

	// Create context with timeout for command execution
	cmdCtx, cancel := context.WithTimeout(*appContext.MainCtx, appContext.Config.MinecraftServer.ServerCommandExecTimeout)
	defer cancel()

	// Execute the command through the MinecraftServer interface
	result, err := (*appContext.MinecraftServer).ExecuteCommand(cmdCtx, message)
	if err != nil {
		slog.Error("Error executing command", "command", message, "error", err.Error())
		errorMsg := models.NewAdmineMessage([]string{"notification"}, "Failed to execute command: "+err.Error())
		eh.pubsub.Publish(appContext.Config.PubSub.AdmineChannelsMap.ServerChannel, errorMsg)
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
	eh.pubsub.Publish(appContext.Config.PubSub.AdmineChannelsMap.ServerChannel, successMsg)

	slog.Info("Executed command successfully", "command", message)
}
