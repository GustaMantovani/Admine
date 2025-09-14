package mcserver

import (
	"fmt"

	"admine.com/server_handler/internal"
	mcserver "admine.com/server_handler/internal/core/mc_server/mc_server_impls"
	"admine.com/server_handler/pkg"
)

func CreateMinecraftServer(serverType string) (MinecraftServer, error) {
	ctx := internal.Get() // singleton AppContext
	if ctx == nil {
		return nil, fmt.Errorf("AppContext not initialized")
	}

	switch serverType {
	case "docker":
		ctx := internal.Get()
		dc := pkg.NewDockerCompose(ctx.Config.DockerComposePath)
		return mcserver.NewDockerMinecraftServer(dc, "mine_server"), nil
	default:
		return nil, fmt.Errorf("unknown pubsub type: %s", serverType)
	}
}
