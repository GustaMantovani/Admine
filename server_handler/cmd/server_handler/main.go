package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/GustaMantovani/Admine/server_handler/internal/api"
	"github.com/GustaMantovani/Admine/server_handler/internal/config"
	"github.com/GustaMantovani/Admine/server_handler/internal/logger"
	"github.com/GustaMantovani/Admine/server_handler/internal/pubsub"
	"github.com/GustaMantovani/Admine/server_handler/internal/server"
)

func main() {
	mainCtx, cancel := context.WithCancel(context.Background())

	configPath := "server_handler_config.yaml"
	if args := os.Args; len(args) > 1 {
		configPath = args[1]
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if err := logger.Setup(cfg.App.LogFilePath, cfg.App.LogLevel); err != nil {
		log.Fatalf("Failed to setup logger: %v", err)
	}

	slog.Info("Server Handler starting...")

	srv, err := server.NewDocker(cfg.MinecraftServer)
	if err != nil {
		slog.Error("Failed to create Minecraft server", "error", err)
		os.Exit(1)
	}

	ps, err := pubsub.NewRedis(cfg.PubSub, mainCtx)
	if err != nil {
		slog.Error("Failed to create PubSub service", "error", err)
		os.Exit(1)
	}

	origin := cfg.App.SelfOriginName
	serverChannel := cfg.PubSub.AdmineChannelsMap.ServerChannel

	router := api.SetupRouter(srv, ps, origin, serverChannel, cfg.App.LogLevel, cfg.MinecraftServer, mainCtx)
	webServer := api.NewServer(cfg.WebSever, router)
	if err := webServer.StartBackground(); err != nil {
		slog.Error("Failed to start web server", "error", err)
		os.Exit(1)
	}

	eventHandler := pubsub.NewEventHandler(srv, ps, origin, serverChannel, cfg.MinecraftServer, mainCtx)

	msgChannel, err := ps.Subscribe(cfg.PubSub.AdmineChannelsMap.CommandChannel)
	if err != nil {
		slog.Error("Failed to subscribe to commands", "error", err)
		os.Exit(1)
	}

	slog.Info("Server Handler started successfully. Listening for messages on channel", "channel", cfg.PubSub.AdmineChannelsMap.CommandChannel)

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		slog.Info("Received shutdown signal. Shutting down gracefully...")
		cancel()
	}()

	for {
		select {
		case <-mainCtx.Done():
			slog.Info("Shutting down...")

			shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer shutdownCancel()

			if err := webServer.Stop(shutdownCtx); err != nil {
				slog.Error("Error stopping web server", "error", err)
			}

			if err := ps.Close(); err != nil {
				slog.Error("Error closing PubSub service", "error", err)
			}

			slog.Info("Server Handler stopped")
			return

		case msg := <-msgChannel:
			if msg != nil {
				slog.Info("Received message with tags", "tags", msg.Tags)
				if err := eventHandler.ManageCommand(msg); err != nil {
					slog.Error("Error processing message", "error", err)
				}
			}
		}
	}
}
