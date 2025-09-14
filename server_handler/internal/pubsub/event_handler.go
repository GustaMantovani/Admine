package pubsub

import (
	"log"
	"strings"

	"admine.com/server_handler/internal"
	"admine.com/server_handler/internal/pubsub/models"
	"admine.com/server_handler/pkg"
)

func ManageCommand(msg *models.AdmineMessage, ps PubSubService) error {
	ctx := internal.Get()

	if msg.HasTag("server_up") {
		serverUp(ps, ctx)
	} else if msg.HasTag("server_down") {
		serverDown(ps, ctx)
	} else if msg.HasTag("command") {
		command(ps, ctx, msg.Message)
	} else {
		responseMsg := models.NewAdmineMessage([]string{"error"}, "Invalid tag.")
		ps.Publish("admine_responses", responseMsg)

		if pkg.Logger != nil {
			pkg.Logger.Error("Received an invalid tag: %v", msg.Tags)
		} else {
			log.Println("Received an invalid tag:", msg.Tags)
		}
	}

	return nil
}

func serverUp(ps PubSubService, ctx *internal.AppContext) {
	// Start the server using the MinecraftServer interface
	if ctx.MinecraftServer == nil {
		pkg.Logger.Error("MinecraftServer is not initialized")
		responseMsg := models.NewAdmineMessage([]string{"error"}, "Server not initialized")
		ps.Publish("admine_responses", responseMsg)
		return
	}

	// Send starting message
	responseMsg := models.NewAdmineMessage([]string{"server_status"}, "Starting server")
	ps.Publish("admine_responses", responseMsg)

	err := (*ctx.MinecraftServer).Start()
	if err != nil {
		if pkg.Logger != nil {
			pkg.Logger.Error("Error starting server: %s", err.Error())
		}
		errorMsg := models.NewAdmineMessage([]string{"error"}, "Failed to start server: "+err.Error())
		ps.Publish("admine_responses", errorMsg)
		return
	}

	// Get server info (this might include network information)
	info, err := (*ctx.MinecraftServer).Info()
	if err != nil {
		if pkg.Logger != nil {
			pkg.Logger.Error("Error getting server info: %s", err.Error())
		}
		info = "Server started successfully"
	}

	successMsg := models.NewAdmineMessage([]string{"server_up"}, info)
	ps.Publish("admine_responses", successMsg)

	if pkg.Logger != nil {
		pkg.Logger.Info("Server started successfully")
	}
}

func serverDown(ps PubSubService, ctx *internal.AppContext) {
	if ctx.MinecraftServer == nil {
		pkg.Logger.Error("MinecraftServer is not initialized")
		responseMsg := models.NewAdmineMessage([]string{"error"}, "Server not initialized")
		ps.Publish("admine_responses", responseMsg)
		return
	}

	// Send stopping message
	responseMsg := models.NewAdmineMessage([]string{"server_status"}, "Stopping server")
	ps.Publish("admine_responses", responseMsg)

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
		ps.Publish("admine_responses", errorMsg)
		return
	}

	successMsg := models.NewAdmineMessage([]string{"server_down"}, "Server stopped successfully")
	ps.Publish("admine_responses", successMsg)

	if pkg.Logger != nil {
		pkg.Logger.Info("Server stopped successfully")
	}
}

func command(ps PubSubService, ctx *internal.AppContext, message string) {
	if ctx.MinecraftServer == nil {
		pkg.Logger.Error("MinecraftServer is not initialized")
		responseMsg := models.NewAdmineMessage([]string{"error"}, "Server not initialized")
		ps.Publish("admine_responses", responseMsg)
		return
	}

	// Execute the command through the MinecraftServer interface
	result, err := (*ctx.MinecraftServer).ExecuteCommand(message)
	if err != nil {
		if pkg.Logger != nil {
			pkg.Logger.Error("Error executing command '%s': %s", message, err.Error())
		}
		errorMsg := models.NewAdmineMessage([]string{"error"}, "Failed to execute command: "+err.Error())
		ps.Publish("admine_responses", errorMsg)
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
	ps.Publish("admine_responses", successMsg)

	if pkg.Logger != nil {
		pkg.Logger.Info("Executed command '%s' successfully", message)
	}
}
