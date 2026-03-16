package pubsub

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/GustaMantovani/Admine/server_handler/internal/config"
	"github.com/GustaMantovani/Admine/server_handler/internal/server"
)

// EventHandler routes incoming pub/sub messages to the appropriate server operations
type EventHandler struct {
	server        server.MinecraftServer
	pubsub        PubSubService
	origin        string
	serverChannel string
	cfg           config.MinecraftServerConfig
	mainCtx       context.Context
}

// NewEventHandler creates a new EventHandler
func NewEventHandler(
	srv server.MinecraftServer,
	ps PubSubService,
	origin string,
	serverChannel string,
	cfg config.MinecraftServerConfig,
	mainCtx context.Context,
) *EventHandler {
	return &EventHandler{
		server:        srv,
		pubsub:        ps,
		origin:        origin,
		serverChannel: serverChannel,
		cfg:           cfg,
		mainCtx:       mainCtx,
	}
}

func (eh *EventHandler) publish(tags []string, message string) {
	msg := NewAdmineMessage(eh.origin, tags, message)
	eh.pubsub.Publish(eh.serverChannel, msg)
}

// ManageCommand processes an incoming message and routes it to the appropriate handler
func (eh *EventHandler) ManageCommand(msg *AdmineMessage) error {
	if eh.server == nil {
		slog.Error("MinecraftServer is not initialized")
		eh.publish([]string{"notification"}, "Server not initialized")
		return fmt.Errorf("minecraft server is not initialized")
	}

	switch {
	case msg.HasTag("server_on"):
		eh.serverUp()
	case msg.HasTag("server_off"):
		eh.serverOff()
	case msg.HasTag("server_down"):
		eh.serverDown()
	case msg.HasTag("restart"):
		eh.restart()
	case msg.HasTag("command"):
		eh.command(msg.Message)
	default:
		eh.publish([]string{"notification"}, "Invalid tag.")
		slog.Error("Received an invalid tag", "tags", msg.Tags)
	}

	return nil
}

func (eh *EventHandler) serverUp() {
	eh.publish([]string{"notification"}, "Starting server")

	startCtx, cancel := context.WithTimeout(eh.mainCtx, eh.cfg.ServerOnTimeout)
	defer cancel()

	if err := eh.server.Start(startCtx); err != nil {
		slog.Error("Error starting server", "error", err.Error())
		eh.publish([]string{"notification"}, "Failed to start server: "+err.Error())
		return
	}

	infoCtx, infoCancel := context.WithTimeout(eh.mainCtx, eh.cfg.ServerCommandExecTimeout)
	defer infoCancel()

	startInfo := eh.server.StartUpInfo(infoCtx)
	slog.Debug("VPN startup info", "node_key", startInfo)
	eh.publish([]string{"server_on"}, startInfo)

	slog.Info("Server started successfully")
}

func (eh *EventHandler) serverOff() {
	eh.publish([]string{"notification"}, "Stopping server")

	stopCtx, cancel := context.WithTimeout(eh.mainCtx, eh.cfg.ServerOffTimeout)
	defer cancel()

	if err := eh.server.Stop(stopCtx); err != nil {
		slog.Error("Error stopping server", "error", err.Error())
		eh.publish([]string{"notification"}, "Error stopping server: "+err.Error())
		return
	}

	eh.publish([]string{"server_off"}, "Server stopped successfully")
	slog.Info("Server stopped successfully")
}

func (eh *EventHandler) serverDown() {
	eh.publish([]string{"notification"}, "Removing server")

	cmdCtx, cmdCancel := context.WithTimeout(eh.mainCtx, eh.cfg.ServerCommandExecTimeout)
	defer cmdCancel()

	if _, err := eh.server.ExecuteCommand(cmdCtx, "/stop"); err != nil {
		slog.Error("Error executing stop command", "error", err.Error())
	}

	downCtx, downCancel := context.WithTimeout(eh.mainCtx, eh.cfg.ServerOffTimeout)
	defer downCancel()

	if err := eh.server.Down(downCtx); err != nil {
		slog.Error("Error stopping server", "error", err.Error())
		eh.publish([]string{"notification"}, "Error removing server: "+err.Error())
		return
	}

	eh.publish([]string{"server_off"}, "Server removed successfully")
	slog.Info("Server removed successfully")
}

func (eh *EventHandler) restart() {
	eh.publish([]string{"notification"}, "Restarting server")
	slog.Info("Starting server restart process")

	stopCtx, stopCancel := context.WithTimeout(eh.mainCtx, eh.cfg.ServerOffTimeout)
	defer stopCancel()

	if err := eh.server.Stop(stopCtx); err != nil {
		slog.Error("Error stopping server for restart", "error", err.Error())
		eh.publish([]string{"notification"}, "Failed to restart server: "+err.Error())
		return
	}

	startCtx, startCancel := context.WithTimeout(eh.mainCtx, eh.cfg.ServerOnTimeout)
	defer startCancel()

	if err := eh.server.Start(startCtx); err != nil {
		slog.Error("Error starting server after restart", "error", err.Error())
		eh.publish([]string{"notification"}, "Failed to start server after stop: "+err.Error())
		return
	}

	infoCtx, infoCancel := context.WithTimeout(eh.mainCtx, eh.cfg.ServerCommandExecTimeout)
	defer infoCancel()

	startInfo := eh.server.StartUpInfo(infoCtx)
	slog.Debug("VPN startup info", "node_key", startInfo)
	eh.publish([]string{"server_on"}, startInfo)

	slog.Info("Server restarted successfully")
}

func (eh *EventHandler) command(message string) {
	cmdCtx, cancel := context.WithTimeout(eh.mainCtx, eh.cfg.ServerCommandExecTimeout)
	defer cancel()

	result, err := eh.server.ExecuteCommand(cmdCtx, message)
	if err != nil {
		slog.Error("Error executing command", "command", message, "error", err.Error())
		eh.publish([]string{"notification"}, "Failed to execute command: "+err.Error())
		return
	}

	responseMessage := "Command executed successfully"
	if strings.TrimSpace(result.Output) != "" {
		responseMessage = result.Output
	}

	eh.publish([]string{"command_result"}, responseMessage)
	slog.Info("Executed command successfully", "command", message)
}
