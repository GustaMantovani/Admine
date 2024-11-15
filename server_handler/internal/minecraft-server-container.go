package internal

import (
	"log"
	"strings"

	docker "github.com/fsouza/go-dockerclient"
)

type MinecraftServerContainerByCompose struct {
	containerName            string
	composeDirectoryFullName string
	containerStatus          string
	dockerClient             *docker.Client
}

func NewMinecraftServerContainerByCompose(client *docker.Client, serviceName, fullNameDirectory string) MinecraftServerContainerByCompose {
	partsDirectory := strings.Split(fullNameDirectory, "/")
	containerName := partsDirectory[len(partsDirectory)-2] + "-" + serviceName + "-1"

	minecraftServer := MinecraftServerContainerByCompose{
		containerName:            containerName,
		composeDirectoryFullName: fullNameDirectory,
		dockerClient:             client,
	}

	return minecraftServer
}

func (ms MinecraftServerContainerByCompose) UpMinecraftServerContainerByCompose() ([]byte, error) {
	return StartServerDockerCompose(ms.composeDirectoryFullName)
}

func (ms *MinecraftServerContainerByCompose) updateMinecraftServerContainerStatus() {
	ms.containerStatus = SeeContainerStatus(ms.dockerClient, ms.containerName)
}

func (ms *MinecraftServerContainerByCompose) VerifyContainerAndUpIfDown() {
	ms.updateMinecraftServerContainerStatus()
	if !strings.Contains(ms.containerStatus, "Up") {
		log.Println("Servidor não está de pé. Status do seu container: ", ms.containerStatus)
		ms.UpMinecraftServerContainerByCompose()
	} else {
		log.Println("Servidor de pé. Status: ", ms.containerStatus)
	}
}
