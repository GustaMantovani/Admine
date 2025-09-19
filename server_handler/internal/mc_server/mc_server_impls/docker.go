package mcserver

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/GustaMantovani/Admine/server_handler/internal/config"
	"github.com/GustaMantovani/Admine/server_handler/pkg"
)

type DockerMinecraftServer struct {
	DockerCompose *pkg.DockerCompose
	DockerConfig  config.DockerConfig
	Context       context.Context
}

func NewDockerMinecraftServer(compose *pkg.DockerCompose, dockerConfig config.DockerConfig, ctx context.Context) *DockerMinecraftServer {
	return &DockerMinecraftServer{
		DockerCompose: compose,
		DockerConfig:  dockerConfig,
		Context:       ctx,
	}
}

func (d *DockerMinecraftServer) Start() error {
	return d.DockerCompose.Up(true)
}

func (d *DockerMinecraftServer) Stop() error {
	if _, err := d.ExecuteCommand("/stop"); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(d.Context, 60*time.Second)
	defer cancel()

	done := make(chan error, 1)
	go func() {
		err := pkg.StreamContainerLogs(ctx, d.DockerConfig.ContainerName, func(line string) {
			slog.Debug("Container line:", "line", line)
			if strings.Contains(line, "All dimensions are saved") {
				done <- nil
			}
		})
		if err != nil {
			done <- err
		}
	}()

	select {
	case <-ctx.Done():
		return fmt.Errorf("timeout esperando shutdown do minecraft")
	case err := <-done:
		if err != nil {
			return err
		}
	}

	// Agora é seguro parar o container
	return d.DockerCompose.Stop()
}

func (d *DockerMinecraftServer) Down() error {
	return d.DockerCompose.Down()
}

func (d *DockerMinecraftServer) Restart() error {
	if err := d.Stop(); err != nil {
		return err
	}
	return d.Start()
}

func (d *DockerMinecraftServer) Status() (string, error) {
	return "nil", nil
}

func (d *DockerMinecraftServer) Info() (string, error) {
	return "nil", nil
}

func (d *DockerMinecraftServer) StartUpInfo() string {
	id, err := pkg.GetZeroTierNodeID(d.DockerConfig.ContainerName)
	if err != nil {
		return ""
	}

	return id
}

func (d *DockerMinecraftServer) ExecuteCommand(command string) (string, error) {
	return "nil", pkg.WriteToContainer(d.Context, d.DockerConfig.ServiceName, command)
}
