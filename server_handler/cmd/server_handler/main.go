package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/GustaMantovani/Admine/server_handler/internal"
	"github.com/GustaMantovani/Admine/server_handler/internal/api"
	"github.com/GustaMantovani/Admine/server_handler/internal/pubsub"
	"github.com/GustaMantovani/Admine/server_handler/pkg"
)

func main() {

	configPath := "server_handler_config.yaml"
	args := os.Args

	if len(args) > 1 {
		configPath = args[1]
	}

	// Initialize application context
	ctx, err := internal.Init(configPath)
	if err != nil {
		log.Fatalf("Failed to initialize app context: %v", err)
	}

	// Initialize slog
	err = pkg.Setup(ctx.Config.App.LogFilePath, ctx.Config.App.LogLevel)
	if err != nil {
		log.Fatalf("Failed to setup slog: %v", err)
	}

	slog.Info("Server Handler starting...")

	// Create PubSub service
	pubsubService, err := pubsub.CreatePubSub(ctx.Config.PubSub)
	if err != nil {
		slog.Error("Failed to create PubSub service", "error", err)
		os.Exit(1)
	}

	// Create event handler
	eventHandler := pubsub.NewEventHandler(pubsubService)

	// Create and start web server in background
	webServer := api.NewServer(ctx.Config)
	go func() {
		if err := webServer.StartBackground(); err != nil {
			slog.Error("Failed to start web server", "error", err)
		}
	}()

	// Subscribe to incoming messages
	msgChannel, err := pubsubService.Subscribe(ctx.Config.PubSub.AdmineChannelsMap.CommandChannel)
	if err != nil {
		slog.Error("Failed to subscribe to commands", "error", err)
		os.Exit(1)
	}

	slog.Info("Server Handler started successfully. Listening for messages on channel", "channel", ctx.Config.PubSub.AdmineChannelsMap.CommandChannel)

	// Create context for graceful shutdown
	mainCtx, cancel := context.WithCancel(context.Background())

	// Handle graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		slog.Info("Received shutdown signal. Shutting down gracefully...")
		cancel()
	}()

	// Main message processing loop
	for {
		select {
		case <-mainCtx.Done():
			slog.Info("Shutting down...")

			// Stop web server
			shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer shutdownCancel()

			if err := webServer.Stop(shutdownCtx); err != nil {
				slog.Error("Error stopping web server", "error", err)
			}

			// Close PubSub connection
			if err := pubsubService.Close(); err != nil {
				slog.Error("Error closing PubSub service", "error", err)
			}

			slog.Info("Server Handler stopped")
			return

		case msg := <-msgChannel:
			if msg != nil {
				slog.Info("Received message with tags", "tags", msg.Tags)

				// Process the message
				if err := eventHandler.ManageCommand(msg); err != nil {
					slog.Error("Error processing message", "error", err)
				}
			}
		}
	}
}
