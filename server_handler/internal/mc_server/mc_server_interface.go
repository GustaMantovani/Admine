package mcserver

import (
	"context"

	"github.com/GustaMantovani/Admine/server_handler/internal/mc_server/models"
)

type MinecraftServer interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Down(ctx context.Context) error
	Restart(ctx context.Context) error
	Status(ctx context.Context) (*models.ServerStatus, error)
	Info(ctx context.Context) (*models.ServerInfo, error)
	StartUpInfo(ctx context.Context) string
	ExecuteCommand(ctx context.Context, command string) (*models.CommandResult, error)
}
