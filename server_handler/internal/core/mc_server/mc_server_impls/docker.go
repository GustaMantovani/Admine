package mcserver

import (
	"admine.com/server_handler/pkg"
)

type DockerMinecraftServer struct {
	DockerCompose *pkg.DockerCompose
	ContainerName string
}

func NewDockerMinecraftServer(compose *pkg.DockerCompose, containerName string) *DockerMinecraftServer {
	return &DockerMinecraftServer{
		DockerCompose: compose,
		ContainerName: containerName,
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

func (d *DockerMinecraftServer) ExecuteCommand(command string) (string, error) {
	return "nil", pkg.WriteToContainer(nil, d.ContainerName, command)
}
