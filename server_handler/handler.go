package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"

	docker "github.com/fsouza/go-dockerclient"
)

func main() {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		log.Fatal(err)
	}

	cmd := exec.Command("docker", "compose", "up", "-d")
	cmd.Dir = "/home/andre/pgm/pessoal/minecraft-server-on-docker/"

	mineNameServiceCompose := "mine_server-1"
	mineServerComposeDirectory := "minecraft-server-on-docker-"
	containerName := mineServerComposeDirectory + mineNameServiceCompose

	var containerStatus string

	for {
		strings.Contains(containerStatus, "Up")

		containerStatus = verificarStatusContainer(client, containerName)

		fmt.Printf("Status do container do servidor: '%s'\n", containerStatus)

		if !strings.Contains(containerStatus, "Up") {
			fmt.Printf("Servidor derrubado.\n")
			cmd := exec.Command("docker", "compose", "up", "-d")
			cmd.Dir = "/home/andre/pgm/pessoal/minecraft-server-on-docker/"
			output, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Println("Erro", err)
			} else {
				fmt.Println("Comando executado com sucesso: ", string(output))
			}

		}
		time.Sleep(1 * time.Second)
	}
}

func verificarStatusContainer(client *docker.Client, containerName string) string {
	var containerStatus string

	containers, err := client.ListContainers(docker.ListContainersOptions{All: true})
	if err != nil {
		log.Fatal(err)
	}

	for _, container := range containers {
		for _, name := range container.Names {
			if name == fmt.Sprintf("/%s", containerName) {
				containerStatus = container.Status
				break
			}
		}
	}

	return containerStatus
}
