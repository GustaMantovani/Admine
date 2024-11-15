package main

import (
	"fmt"
	"log"
	"os/exec"
	"server/handler/internal"
	"strings"
	"time"

	docker "github.com/fsouza/go-dockerclient"
)

func main() {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		log.Fatal(err)
	}

	mineNameServiceCompose := "mine_server-1"
	mineServerComposeDirectory := "minecraft-server-on-docker-"
	containerName := mineServerComposeDirectory + mineNameServiceCompose

	composeDirectory := "/home/andre/pgm/pessoal/minecraft-server-on-docker/"
	internal.StartServerDockerCompose(composeDirectory)

	var containerStatus string

	for {
		strings.Contains(containerStatus, "Up")

		containerStatus = internal.VerificarStatusContainer(client, containerName)

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
