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
	ComposeDirectoryFullName string
	containerId              string
	client                   *dockerClient.Client
}

func NewMinecraftServerContainerByCompose() MinecraftServerContainerByCompose {
	// partsDirectory := strings.Split(fullNameDirectory, "/")
	// containerName := partsDirectory[len(partsDirectory)-1] + "-" + serviceName + "-1"

	client, err := dockerClient.NewClientFromEnv()
	if err != nil {
		log.Fatal("Erro relacionado ao cliente docker: ", err)
	}

	minecraftServer := MinecraftServerContainerByCompose{
		ContainerName:            "",
		ComposeDirectoryFullName: "",
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
	ms.ComposeDirectoryFullName = fullNameDirectory
}

func (ms MinecraftServerContainerByCompose) UpMinecraftServerContainerByCompose() ([]byte, error) {
	return docker.StartServerDockerCompose(ms.ComposeDirectoryFullName)
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

// Pega as informações do servidor por argumentos
func (ms *MinecraftServerContainerByCompose) ConfigureWithArgs(args []string) {
	ms.SetContainerNameByServiceAndDirectory(args[0], file.GetLocalDirectory())
}

// Pega as informações do servidor por variáveis de ambiente
func (ms *MinecraftServerContainerByCompose) ConfigureWithEnv() {
	serverName := os.Getenv("MINECRAFT_SERVER_SERVICE")
	directory := os.Getenv("MINECRAFT_SERVER_DIRECTORY")

	ms.SetContainerNameByServiceAndDirectory(serverName, directory)
}

// Pega as informações do servidor pelo arquivo de configuração
func (ms *MinecraftServerContainerByCompose) ConfigureWithFile() {
	configFileData, err := file.GetConfigFileData()
	if err != nil {
		fmt.Println("Erro: ", err)
		os.Exit(1)
	}

	serverName := configFileData.ServerName
	directory := configFileData.ComposeDirectory

	ms.SetContainerNameByServiceAndDirectory(serverName, directory)
}

/*
Configura os metadados do servidor, verificando se vão ser usados argumentos,
variáveis de ambiente ou um arquivo de configuração para tal.

Se os parâmetros 'env' e 'file' forem true, significa que duas flags excludentes foram chamadas.

Se houver argumentos, estes serão usados para configurar o servidor, indepente das flags.
*/
func (ms *MinecraftServerContainerByCompose) ConfigureMinecraftServer(env, file bool, args []string) {
	if len(args) != 0 {
		ms.ConfigureWithArgs(args)
	} else {
		if env && file {
			log.Fatal("Flags excludentes foram chamadas.")
		} else if env {
			ms.ConfigureWithEnv()
		} else if file {
			ms.ConfigureWithFile()
		}
	}
}
