package mcserver

import (
	"fmt"

	"github.com/GustaMantovani/Admine/server_handler/internal/config"
	mcserver "github.com/GustaMantovani/Admine/server_handler/internal/mc_server/mc_server_impls"
	"github.com/GustaMantovani/Admine/server_handler/pkg"
)

func CreateMinecraftServer(cfg *config.Config) (MinecraftServer, error) {
	switch cfg.MinecraftServer.RuntimeType {
	case "docker":
		dc := pkg.NewDockerCompose(cfg.MinecraftServer.Docker.ComposeOutputPath)
		return mcserver.NewDockerMinecraftServer(dc, cfg), nil
	default:
		return nil, fmt.Errorf("unknown runtime type: %s", cfg.MinecraftServer.RuntimeType)
	}
}
