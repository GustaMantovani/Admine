package internal

import (
	"fmt"
	"log"
	"os"
	"strings"

	docker "github.com/fsouza/go-dockerclient"
)

type MinecraftServerContainerByCompose struct {
	ContainerName            string
	composeDirectoryFullName string
	containerStatus          string
	DockerClient             *docker.Client
}

func NewMinecraftServerContainerByCompose(client *docker.Client, serviceName, fullNameDirectory string) MinecraftServerContainerByCompose {
	partsDirectory := strings.Split(fullNameDirectory, "/")
	containerName := partsDirectory[len(partsDirectory)-2] + "-" + serviceName + "-1"

	minecraftServer := MinecraftServerContainerByCompose{
		ContainerName:            containerName,
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

	ms.ContainerName = containerName
	ms.composeDirectoryFullName = fullNameDirectory
}

func (ms MinecraftServerContainerByCompose) UpMinecraftServerContainerByCompose() ([]byte, error) {
	return StartServerDockerCompose(ms.composeDirectoryFullName)
}

func (ms *MinecraftServerContainerByCompose) updateMinecraftServerContainerStatus() {
	ms.containerStatus = SeeContainerStatus(ms.DockerClient, ms.ContainerName)
}

// Verifica se o container do servidor está de pé e se não estiver sobe ele
func (ms *MinecraftServerContainerByCompose) VerifyContainerAndUpIfDown() (string, bool) {
	ms.updateMinecraftServerContainerStatus()
	var msg string
	if !strings.Contains(ms.containerStatus, "Up") {
		msg = "Servidor não está de pé. Status do seu container: " + ms.containerStatus
		log.Println(msg)
		ms.UpMinecraftServerContainerByCompose()
		return msg, false
	} else {
		msg = "Servidor de pé. Status: " + ms.containerStatus
		log.Println(msg)
		return msg, true
	}
}

// Pega as informações do servidor de argumentos
func (ms *MinecraftServerContainerByCompose) ConfigureWithArgs(args []string) {
	ms.SetContainerNameByServiceAndDirectory(args[0], getLocalDirectory())
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
	configFileData, err := GetConfigFileData()
	if err != nil {
		fmt.Println("Erro: ", err)
		os.Exit(1)
	}

	serverName := configFileData.ServerName
	directory := configFileData.ComposeDirectory

	ms.SetContainerNameByServiceAndDirectory(serverName, directory)
	ms.UpMinecraftServerContainerByCompose()
}
