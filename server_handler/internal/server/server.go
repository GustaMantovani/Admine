package server

import (
	"context"
	"fmt"
	"io"

	"github.com/GustaMantovani/Admine/server_handler/internal/config"
	"github.com/GustaMantovani/Admine/server_handler/internal/docker"
)

// MinecraftServer defines the operations supported on a Minecraft server instance
type MinecraftServer interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Down(ctx context.Context) error
	Restart(ctx context.Context) error
	Status(ctx context.Context) (*ServerStatus, error)
	Info(ctx context.Context) (*ServerInfo, error)
	Logs(ctx context.Context, n int) ([]string, error)
	StartUpInfo(ctx context.Context) string
	ExecuteCommand(ctx context.Context, command string) (*CommandResult, error)
	InstallMod(ctx context.Context, fileName string, modData io.Reader) (*ModInstallResult, error)
	ListMods(ctx context.Context) (*ModListResult, error)
	RemoveMod(ctx context.Context, fileName string) (*ModInstallResult, error)
}

// NewDocker creates a MinecraftServer backed by Docker Compose
func NewDocker(cfg config.MinecraftServerConfig) (MinecraftServer, error) {
	switch cfg.RuntimeType {
	case "docker":
		dc := docker.NewDockerCompose(cfg.Docker.ComposeOutputPath)
		return newDockerMinecraftServer(dc, cfg), nil
	default:
		return nil, fmt.Errorf("unknown runtime type: %s", cfg.RuntimeType)
	}
}
