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
}

func NewDockerMinecraftServer(compose *pkg.DockerCompose, dockerConfig config.DockerConfig) *DockerMinecraftServer {
	return &DockerMinecraftServer{
		DockerCompose: compose,
		DockerConfig:  dockerConfig,
	}
}

func (d *DockerMinecraftServer) Start(ctx context.Context) error {
	return d.DockerCompose.Up(true)
}

func (d *DockerMinecraftServer) Stop(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	done := make(chan error, 1)

	if _, err := d.ExecuteCommand(ctx, "/stop"); err != nil {
		return err
	}

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

	return d.DockerCompose.Stop()
}

func (d *DockerMinecraftServer) Down(ctx context.Context) error {
	return d.DockerCompose.Down()
}

func (d *DockerMinecraftServer) Restart(ctx context.Context) error {
	if err := d.Stop(ctx); err != nil {
		return err
	}
	return d.Start(ctx)
}

func (d *DockerMinecraftServer) Status(ctx context.Context) (string, error) {
	return "nil", nil
}

func (d *DockerMinecraftServer) Info(ctx context.Context) (string, error) {
	return "nil", nil
}

func (d *DockerMinecraftServer) StartUpInfo(ctx context.Context) string {
	id, err := pkg.GetZeroTierNodeID(d.DockerConfig.ContainerName)
	if err != nil {
		return ""
	}

	return id
}

func (d *DockerMinecraftServer) ExecuteCommand(ctx context.Context, command string) (string, error) {
	return "nil", pkg.WriteToContainer(ctx, d.DockerConfig.ServiceName, command)
}
