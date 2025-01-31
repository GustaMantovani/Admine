package minecraftserver

import (
	"fmt"
	"log"
	"os"
	"server/handler/internal/docker"
	"server/handler/internal/file"
	"strings"

	dockerClient "github.com/fsouza/go-dockerclient"
)

// Minecraft Server metadata about the compose container
type MinecraftServerContainerByCompose struct {
	ContainerName            string
	composeDirectoryFullName string
	containerId              string
	client                   *dockerClient.Client
}

func NewMinecraftServerContainerByCompose(serviceName, fullNameDirectory string) MinecraftServerContainerByCompose {
	partsDirectory := strings.Split(fullNameDirectory, "/")
	containerName := partsDirectory[len(partsDirectory)-1] + "-" + serviceName + "-1"

	client, err := dockerClient.NewClientFromEnv()
	if err != nil {
		log.Fatal("Erro relacionado ao cliente docker: ", err)
	}

	minecraftServer := MinecraftServerContainerByCompose{
		ContainerName:            containerName,
		composeDirectoryFullName: fullNameDirectory,
		client:                   client,
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

	ms.ContainerName = containerName
	ms.composeDirectoryFullName = fullNameDirectory
}

func (ms MinecraftServerContainerByCompose) UpMinecraftServerContainerByCompose() ([]byte, error) {
	return docker.StartServerDockerCompose(ms.composeDirectoryFullName)
}

func (ms MinecraftServerContainerByCompose) SeeStatus() string {
	containers, err := ms.client.ListContainers(dockerClient.ListContainersOptions{All: true})
	if err != nil {
		log.Fatal(err)
	}

	for _, container := range containers {
		for _, name := range container.Names {
			if name == fmt.Sprintf("/%s", ms.ContainerName) {
				return container.Status
			}
		}
	}

	return "Container não encontrado."

}

// Pega as informações do servidor de argumentos
func (ms *MinecraftServerContainerByCompose) ConfigureWithArgs(args []string) {
	ms.SetContainerNameByServiceAndDirectory(args[0], file.GetLocalDirectory())
}

// Pega as informações do servidor de variáveis de ambiente
func (ms *MinecraftServerContainerByCompose) ConfigureWithEnv() {
	serverName := os.Getenv("MINECRAFT_SERVER_SERVICE")
	directory := os.Getenv("MINECRAFT_SERVER_DIRECTORY")

	ms.SetContainerNameByServiceAndDirectory(serverName, directory)
	ms.UpMinecraftServerContainerByCompose()
}

// Pega as informações do servidor do arquivo de configuração
func (ms *MinecraftServerContainerByCompose) ConfigureWithFile() {
	configFileData, err := file.GetConfigFileData()
	if err != nil {
		fmt.Println("Erro: ", err)
		os.Exit(1)
	}

	serverName := configFileData.ServerName
	directory := configFileData.ComposeDirectory

	ms.SetContainerNameByServiceAndDirectory(serverName, directory)
	ms.UpMinecraftServerContainerByCompose()
}
