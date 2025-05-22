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

	// Docker client
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return "", err
	}
	defer cli.Close()

	containerName := c.ComposeContainerName

	// Search container by name
	containerID := getContainerId(containerName, cli)

	// Catch logs
	out, err := cli.ContainerLogs(ctx, containerID, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: false,
		Timestamps: false,
		Tail:       "1", // Read last line
	})

	if err != nil {
		return "", err
	}

	defer out.Close()

	// Docker logs comes with a 8 bytes header in strem
	// Needs skip this header
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

	// Create Docker client
	cli, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		log.Fatalf("error creating Docker client: %v", err)
	}

	containerName := c.ComposeContainerName

	// Verifys if container exists
	_, err = cli.ContainerInspect(ctx, containerName)
	if err != nil {
		log.Fatalf("Container not found: %v", err)
	}

	err = waitForContainerStart(cli, containerName)
	if err != nil {
		log.Fatalf("Erro: %v", err)
	}

	log.Println("Container is up")
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
		return ""
	}

	return containerID
}

func waitForContainerStart(cli *client.Client, containerName string) error {
	ctx := context.Background()

	for {
		// List container (stopped containers too)
		containers, err := cli.ContainerList(ctx, container.ListOptions{All: true})
		if err != nil {
			return fmt.Errorf("error listing containers: %v", err)
		}

		if len(containers) > 0 {
			container := containers[0]
			if container.State == "running" {
				return nil
			}
			log.Printf("Container status: %s\n", container.State)
		} else {
			log.Println("Container not found, waiting...")
		}

		time.Sleep(1 * time.Second)
	}
}

func VerifyIfContainerExists() bool {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())

	_, err = cli.ContainerInspect(context.Background(), c.ComposeContainerName)
	if err != nil {
		return false
	}

	return true
}
