package docker

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"os/exec"
	"strings"

	"server_handler/internal/config"

	"slices"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func ReadLastContainerLine() (string, error) {
	ctx := context.Background()

	// Cliente Docker
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return "", err
	}
	defer cli.Close()

	containerName := config.GetInstance().ComposeContainerName

	// Procura o container pelo nome
	containers, err := cli.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		return "", err
	}

	var containerID string
	for _, c := range containers {
		if slices.Contains(c.Names, "/"+containerName) {
			containerID = c.ID
		}
	}

	if containerID == "" {
		return "", fmt.Errorf("container '%s' not found", containerName)
	}

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
