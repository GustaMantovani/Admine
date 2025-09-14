package mcserver

import (
	"fmt"

	"admine.com/server_handler/internal/config"
	mcserver "admine.com/server_handler/internal/mc_server/mc_server_impls"
	"admine.com/server_handler/pkg"
)

func CreateMinecraftServer(config config.MinecraftServerConfig) (MinecraftServer, error) {
	switch config.RuntimeType {
	case "docker":
		dc := pkg.NewDockerCompose(config.Docker.ComposePath)
		return mcserver.NewDockerMinecraftServer(dc, config.Docker.ContainerName, config.Docker.ServiceName), nil
	default:
		return nil, fmt.Errorf("unknown pubsub type: %s", config.RuntimeType)
	}
}
