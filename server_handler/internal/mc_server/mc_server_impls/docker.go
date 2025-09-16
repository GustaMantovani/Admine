package mcserver

import (
	"context"

	"github.com/GustaMantovani/Admine/server_handler/pkg"
)

type DockerMinecraftServer struct {
	DockerCompose *pkg.DockerCompose
	ContainerName string
	ServiceName   string
}

func NewDockerMinecraftServer(compose *pkg.DockerCompose, containerName string, serviceName string) *DockerMinecraftServer {
	return &DockerMinecraftServer{
		DockerCompose: compose,
		ContainerName: containerName,
		ServiceName:   serviceName,
	}
}

func (d *DockerMinecraftServer) Start() error {
	return d.DockerCompose.Up(true)
}

func (d *DockerMinecraftServer) Stop() error {
	return d.DockerCompose.Stop()
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
	id, err := pkg.GetZeroTierNodeID(d.ContainerName)
	if err != nil {
		pkg.Logger.Error("Failed to get ZeroTier Node ID: %v", err)
		return ""
	}

	return id
}

func (d *DockerMinecraftServer) ExecuteCommand(command string) (string, error) {
	return "nil", pkg.WriteToContainer(context.Background(), d.ServiceName, command)
}
