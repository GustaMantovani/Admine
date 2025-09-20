package mcserver

import "context"

type MinecraftServer interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Down(ctx context.Context) error
	Restart(ctx context.Context) error
	Status(ctx context.Context) (string, error)
	Info(ctx context.Context) (string, error)
	StartUpInfo(ctx context.Context) string
	ExecuteCommand(ctx context.Context, command string) (string, error)
}
