package mcserver

import (
	"context"
	"log/slog"
	"strings"

	"github.com/GustaMantovani/Admine/server_handler/internal/config"
	"github.com/GustaMantovani/Admine/server_handler/internal/mc_server/models"
	"github.com/GustaMantovani/Admine/server_handler/pkg"
	"github.com/gorcon/rcon"
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
		return ctx.Err()
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

func (d *DockerMinecraftServer) Status(ctx context.Context) (*models.ServerStatus, error) {
	conn, err := rcon.Dial(d.DockerConfig.RconAddress, d.DockerConfig.RconPassword)
	if err != nil {
		return models.NewServerStatus(
			models.HealthUnknown,
			models.StatusOffline,
			"Server is offline - cannot connect via RCON",
			"0h 0m",
			0.0,
		), nil
	}
	defer conn.Close()

	listResponse, err := conn.Execute("list")
	if err != nil {
		return models.NewServerStatus(
			models.HealthCritical,
			models.StatusUnknown,
			"Server RCON error: "+err.Error(),
			"Unknown",
			0.0,
		), nil
	}

	tps := 20.0 // Default
	if tpsResponse, err := conn.Execute("forge tps"); err == nil {
		if strings.Contains(tpsResponse, "TPS") {
			tps = 20.0
		}
	}

	return models.NewServerStatus(
		models.HealthHealthy,
		models.StatusOnline,
		"Server is online - "+listResponse,
		"Unknown",
		tps,
	), nil
}

func (d *DockerMinecraftServer) Info(ctx context.Context) (*models.ServerInfo, error) {
	conn, err := rcon.Dial(d.DockerConfig.RconAddress, d.DockerConfig.RconPassword)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	seed := "Unknown"
	if seedResponse, err := conn.Execute("seed"); err == nil {
		seed = seedResponse
	}

	minecraftVersion := "1.20.1"
	javaVersion := "17.0.2"
	modEngine := "Vanilla"
	maxPlayers := 20

	if listResponse, err := conn.Execute("list"); err == nil {
		if strings.Contains(listResponse, "max of") {
		}
	}

	return models.NewServerInfo(
		minecraftVersion,
		javaVersion,
		modEngine,
		seed,
		maxPlayers,
	), nil
}

func (d *DockerMinecraftServer) StartUpInfo(ctx context.Context) string {
	id, err := pkg.GetZeroTierNodeID(d.DockerConfig.ContainerName)
	if err != nil {
		return ""
	}

	return id
}

func (d *DockerMinecraftServer) ExecuteCommand(ctx context.Context, command string) (*models.CommandResult, error) {
	conn, err := rcon.Dial(d.DockerConfig.RconAddress, d.DockerConfig.RconPassword)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	response, err := conn.Execute(command)
	if err != nil {
		return nil, err
	}

	slog.Debug(response)

	return models.NewCommandResultWithOutput(response), nil
}
