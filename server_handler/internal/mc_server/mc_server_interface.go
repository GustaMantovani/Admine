package mcserver

import (
	"context"
	"io"

	"github.com/GustaMantovani/Admine/server_handler/internal/mc_server/models"
)

type MinecraftServer interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Down(ctx context.Context) error
	Restart(ctx context.Context) error
	Status(ctx context.Context) (*models.ServerStatus, error)
	Info(ctx context.Context) (*models.ServerInfo, error)
	Logs(ctx context.Context, n int) ([]string, error)
	StartUpInfo(ctx context.Context) string
	ExecuteCommand(ctx context.Context, command string) (*models.CommandResult, error)
	InstallMod(ctx context.Context, fileName string, modData io.Reader) (*models.ModInstallResult, error)
	ListMods(ctx context.Context) (*models.ModListResult, error)
	RemoveMod(ctx context.Context, fileName string) (*models.ModInstallResult, error)
}
