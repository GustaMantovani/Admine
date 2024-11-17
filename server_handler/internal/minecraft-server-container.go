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
	DockerClient             *docker.Client
}

func NewMinecraftServerContainerByCompose(client *docker.Client, serviceName, fullNameDirectory string) MinecraftServerContainerByCompose {
	partsDirectory := strings.Split(fullNameDirectory, "/")
	containerName := partsDirectory[len(partsDirectory)-2] + "-" + serviceName + "-1"

	minecraftServer := MinecraftServerContainerByCompose{
		containerName:            containerName,
		composeDirectoryFullName: fullNameDirectory,
		DockerClient:             client,
	}

	return minecraftServer
}

func (ms *MinecraftServerContainerByCompose) SetContainerNameByServiceAndDirectory(serviceName, fullNameDirectory string) {
	partsDirectory := strings.Split(fullNameDirectory, "/")

	var result []string
	for _, str := range partsDirectory {
		if str != "" {
			result = append(result, str)
		}
	}

	containerName := result[len(result)-1] + "-" + serviceName + "-1"

	ms.containerName = containerName
	ms.composeDirectoryFullName = fullNameDirectory
}

func (ms MinecraftServerContainerByCompose) UpMinecraftServerContainerByCompose() ([]byte, error) {
	return StartServerDockerCompose(ms.composeDirectoryFullName)
}

func (ms *MinecraftServerContainerByCompose) updateMinecraftServerContainerStatus() {
	ms.containerStatus = SeeContainerStatus(ms.DockerClient, ms.containerName)
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
