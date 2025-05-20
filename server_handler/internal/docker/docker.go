package docker

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os/exec"
	"strings"
	"time"

	"server_handler/internal/config"

	"slices"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

var c = config.GetInstance()

func ReadLastContainerLine() (string, error) {
	ctx := context.Background()

	// Cliente Docker
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return "", err
	}
	defer cli.Close()

	containerName := c.ComposeContainerName

	// Procura o container pelo nome
	containerID := getContainerId(containerName, cli)

	// Captura os logs
	out, err := cli.ContainerLogs(ctx, containerID, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: false,
		Timestamps: false,
		Tail:       "1", // Lê apenas a última linha
	})

	if err != nil {
		return "", err
	}

	defer out.Close()

	// Docker logs vêm com um cabeçalho de 8 bytes por stream
	// Precisamos pular esse cabeçalho para obter o conteúdo real
	var buf bytes.Buffer
	_, err = io.Copy(&buf, out)
	if err != nil {
		return "", err
	}

	scanner := bufio.NewScanner(&buf)
	if scanner.Scan() {
		return scanner.Text(), nil
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", fmt.Errorf("no lines found")
}

func GetZeroTierNodeID(containerName string) string {
	cmd := exec.Command("docker", "exec", "-i", containerName, "/bin/bash", "-c", "zerotier-cli info")
	output, err := cmd.CombinedOutput()

	if err != nil {
		panic(err)
	}

	outputStr := string(output)

	parts := strings.Split(outputStr, " ")

	return parts[2]
}

func WaitForBuildAndStart() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Criar cliente Docker
	cli, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		log.Fatalf("Erro ao criar cliente Docker: %v", err)
	}

	containerName := c.ComposeContainerName // Substitua pelo nome do seu container

	// Verificar se o container existe
	_, err = cli.ContainerInspect(ctx, containerName)
	if err != nil {
		log.Fatalf("Container não encontrado: %v", err)
	}

	err = waitForContainerStart(cli, containerName)
	if err != nil {
		log.Fatalf("Erro: %v", err)
	}

	log.Println("Container está rodando com sucesso!")
}

func getContainerId(containerName string, cli *client.Client) string {

	containers, err := cli.ContainerList(context.Background(), container.ListOptions{All: true})
	if err != nil {
		return ""
	}

	var containerID string
	for _, c := range containers {
		if slices.Contains(c.Names, "/"+containerName) {
			containerID = c.ID
		}
	}

	if containerID == "" {
		// return "", fmt.Errorf("container '%s' not found", containerName)
		return ""
	}

	return containerID
}

func waitForContainerStart(cli *client.Client, containerName string) error {
	ctx := context.Background()
	// filter := filters.NewArgs(filters.Arg("name", containerName))

	for {
		// Listar containers (incluindo os que não estão running)
		containers, err := cli.ContainerList(ctx, container.ListOptions{All: true})
		if err != nil {
			return fmt.Errorf("erro ao listar containers: %v", err)
		}

		if len(containers) > 0 {
			container := containers[0]
			if container.State == "running" {
				return nil
			}
			log.Printf("Container status: %s\n", container.State)
		} else {
			log.Println("Container não encontrado, aguardando...")
		}

		time.Sleep(1 * time.Second)
	}
}
