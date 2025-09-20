package mcserver

import (
	"fmt"

	"github.com/GustaMantovani/Admine/server_handler/internal/config"
	mcserver "github.com/GustaMantovani/Admine/server_handler/internal/mc_server/mc_server_impls"
	"github.com/GustaMantovani/Admine/server_handler/pkg"
)

func CreateMinecraftServer(config config.MinecraftServerConfig) (MinecraftServer, error) {
	switch config.RuntimeType {
	case "docker":
		dc := pkg.NewDockerCompose(config.Docker.ComposePath)
		return mcserver.NewDockerMinecraftServer(dc, config.Docker), nil
	default:
		return nil, fmt.Errorf("unknown pubsub type: %s", config.RuntimeType)
	}
}
