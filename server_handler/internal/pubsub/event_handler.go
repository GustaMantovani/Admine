package pubsub

import (
	"log"
	"strings"

	"admine.com/server_handler/internal"
	"admine.com/server_handler/internal/pubsub/models"
	"admine.com/server_handler/pkg"
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
	} else if msg.HasTag("server_down") {
		eh.serverDown()
	} else if msg.HasTag("command") {
		eh.command(msg.Message)
	} else {
		ctx := internal.Get()
		responseMsg := models.NewAdmineMessage([]string{"error"}, "Invalid tag.")
		eh.pubsub.Publish(ctx.Config.PubSub.AdmineChannelsMap.ServerChannel, responseMsg)

		if pkg.Logger != nil {
			pkg.Logger.Error("Received an invalid tag: %v", msg.Tags)
		} else {
			log.Println("Received an invalid tag:", msg.Tags)
		}
	}

	return nil
}

func (eh *EventHandler) serverUp() {
	ctx := internal.Get()
	// Start the server using the MinecraftServer interface
	if ctx.MinecraftServer == nil {
		pkg.Logger.Error("MinecraftServer is not initialized")
		responseMsg := models.NewAdmineMessage([]string{"error"}, "Server not initialized")
		eh.pubsub.Publish(ctx.Config.PubSub.AdmineChannelsMap.ServerChannel, responseMsg)
		return
	}

	// Send starting message

	err := (*ctx.MinecraftServer).Start()
	if err != nil {
		pkg.Logger.Error("Error starting server: %s", err.Error())
		errorMsg := models.NewAdmineMessage([]string{"error"}, "Failed to start server: "+err.Error())
		eh.pubsub.Publish(ctx.Config.PubSub.AdmineChannelsMap.ServerChannel, errorMsg)
		return
	}

	// Get server startInfo (this might include network information)
	startInfo := (*ctx.MinecraftServer).StartUpInfo()

	successMsg := models.NewAdmineMessage([]string{"server_on"}, startInfo)
	eh.pubsub.Publish(ctx.Config.PubSub.AdmineChannelsMap.ServerChannel, successMsg)

	pkg.Logger.Info("Server started successfully")
}

func (eh *EventHandler) serverDown() {
	ctx := internal.Get()
	if ctx.MinecraftServer == nil {
		pkg.Logger.Error("MinecraftServer is not initialized")
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
		if pkg.Logger != nil {
			pkg.Logger.Error("Error executing stop command: %s", err.Error())
		}
	}

	// Wait for graceful shutdown (simplified approach)
	// In a real implementation, you might want to monitor server logs
	// or implement a more sophisticated shutdown detection

	err = (*ctx.MinecraftServer).Stop()
	if err != nil {
		if pkg.Logger != nil {
			pkg.Logger.Error("Error stopping server: %s", err.Error())
		}
		errorMsg := models.NewAdmineMessage([]string{"error"}, "Error stopping server: "+err.Error())
		eh.pubsub.Publish(ctx.Config.PubSub.AdmineChannelsMap.ServerChannel, errorMsg)
		return
	}

	successMsg := models.NewAdmineMessage([]string{"server_down"}, "Server stopped successfully")
	eh.pubsub.Publish(ctx.Config.PubSub.AdmineChannelsMap.ServerChannel, successMsg)

	if pkg.Logger != nil {
		pkg.Logger.Info("Server stopped successfully")
	}
}

func (eh *EventHandler) command(message string) {
	ctx := internal.Get()
	if ctx.MinecraftServer == nil {
		pkg.Logger.Error("MinecraftServer is not initialized")
		responseMsg := models.NewAdmineMessage([]string{"error"}, "Server not initialized")
		eh.pubsub.Publish(ctx.Config.PubSub.AdmineChannelsMap.ServerChannel, responseMsg)
		return
	}

	// Execute the command through the MinecraftServer interface
	result, err := (*ctx.MinecraftServer).ExecuteCommand(message)
	if err != nil {
		if pkg.Logger != nil {
			pkg.Logger.Error("Error executing command '%s': %s", message, err.Error())
		}
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

	if pkg.Logger != nil {
		pkg.Logger.Info("Executed command '%s' successfully", message)
	}
}
