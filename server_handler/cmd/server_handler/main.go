package main

import (
	"context"
	"log"
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

	// Initialize logger
	logger, err := pkg.Setup(ctx.Config.App.LogFilePath)
	if err != nil {
		log.Fatalf("Failed to setup logger: %v", err)
	}

	logger.Info("Server Handler starting...")

	// Create PubSub service
	pubsubService, err := pubsub.CreatePubSub(ctx.Config.PubSub)
	if err != nil {
		logger.Error("Failed to create PubSub service: %v", err)
		os.Exit(1)
	}

	// Create event handler
	eventHandler := pubsub.NewEventHandler(pubsubService)

	// Create and start web server in background
	webServer := api.NewServer(ctx.Config)
	go func() {
		if err := webServer.StartBackground(); err != nil {
			logger.Error("Failed to start web server: %v", err)
		}
	}()

	// Subscribe to incoming messages
	msgChannel, err := pubsubService.Subscribe(ctx.Config.PubSub.AdmineChannelsMap.CommandChannel)
	if err != nil {
		logger.Error("Failed to subscribe to commands: %v", err)
		os.Exit(1)
	}

	logger.Info("Server Handler started successfully. Listening for messages on channel: %s", ctx.Config.PubSub.AdmineChannelsMap.CommandChannel)

	// Create context for graceful shutdown
	mainCtx, cancel := context.WithCancel(context.Background())

	// Handle graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		logger.Info("Received shutdown signal. Shutting down gracefully...")
		cancel()
	}()

	// Main message processing loop
	for {
		select {
		case <-mainCtx.Done():
			logger.Info("Shutting down...")

			// Stop web server
			shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer shutdownCancel()

			if err := webServer.Stop(shutdownCtx); err != nil {
				logger.Error("Error stopping web server: %v", err)
			}

			// Close PubSub connection
			if err := pubsubService.Close(); err != nil {
				logger.Error("Error closing PubSub service: %v", err)
			}

			logger.Info("Server Handler stopped")
			return

		case msg := <-msgChannel:
			if msg != nil {
				logger.Info("Received message with tags: %v", msg.Tags)

				// Process the message
				if err := eventHandler.ManageCommand(msg); err != nil {
					logger.Error("Error processing message: %v", err)
				}
			}
		}
	}
}
