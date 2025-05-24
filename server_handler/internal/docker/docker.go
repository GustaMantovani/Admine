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
		log.Printf("Failed to create Docker client: %v", err)
		return "", err
	}
	defer cli.Close()

	containerName := c.ComposeContainerName
	log.Printf("Reading last container line from: %s", containerName)

	// Search container by name
	containerID := getContainerId(containerName, cli)
	if containerID == "" {
		log.Printf("Container %s not found", containerName)
		return "", fmt.Errorf("container not found: %s", containerName)
	}

	// Get logs
	log.Printf("Fetching logs from container: %s", containerID)
	out, err := cli.ContainerLogs(ctx, containerID, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: false,
		Timestamps: false,
		Tail:       "1", // Read last line
	})

	if err != nil {
		log.Printf("Failed to get container logs: %v", err)
		return "", err
	}

	defer out.Close()

	// Docker logs come with an 8-byte header in stream
	// Need to skip this header
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
	log.Printf("Getting ZeroTier Node ID from container: %s", containerName)

	cmd := exec.Command("docker", "exec", "-i", containerName, "/bin/bash", "-c", "zerotier-cli info")
	output, err := cmd.CombinedOutput()

	if err != nil {
		log.Printf("Failed to get ZeroTier Node ID: %v", err)
		panic(err)
	}

	outputStr := string(output)
	log.Printf("ZeroTier CLI output: %s", outputStr)

	parts := strings.Split(outputStr, " ")
	if len(parts) < 3 {
		log.Printf("Unexpected ZeroTier CLI output format: %s", outputStr)
		panic("Invalid ZeroTier CLI output format")
	}

	nodeID := parts[2]
	log.Printf("Retrieved ZeroTier Node ID: %s", nodeID)
	return nodeID
}

func WaitForBuildAndStart() {
	log.Println("Waiting for container to build and start...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Create Docker client
	cli, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		log.Fatalf("Error creating Docker client: %v", err)
	}

	containerName := c.ComposeContainerName
	log.Printf("Waiting for container: %s", containerName)

	// Verifys if container exists
	_, err = cli.ContainerInspect(ctx, containerName)
	if err != nil {
		log.Fatalf("Container not found after start: %v", err)
	}

	err = waitForContainerStart(cli, containerName)
	if err != nil {
		log.Fatalf("Erro: %v", err)
	}

	log.Println("Container is up")
}

func getContainerId(containerName string, cli *client.Client) string {
	log.Printf("Searching for container with name: %s", containerName)

	containers, err := cli.ContainerList(context.Background(), container.ListOptions{All: true})
	if err != nil {
		log.Printf("Failed to list containers: %v", err)
		return ""
	}

	var containerID string
	for _, c := range containers {
		if slices.Contains(c.Names, "/"+containerName) {
			containerID = c.ID
			log.Printf("Found container %s with ID: %s", containerName, containerID)
			break
		}
	}

	if containerID == "" {
		log.Printf("Container %s not found", containerName)
		return ""
	}

	return containerID
}

func waitForContainerStart(cli *client.Client, containerName string) error {
	ctx := context.Background()
	log.Printf("Waiting for container %s to start...", containerName)

	for {
		// List containers (including stopped ones)
		containers, err := cli.ContainerList(ctx, container.ListOptions{All: true})
		if err != nil {
			log.Printf("Error listing containers: %v", err)
			return fmt.Errorf("error listing containers: %v", err)
		}

		if len(containers) > 0 {
			container := containers[0]
			log.Printf("Container status: %s", container.State)
			if container.State == "running" {
				log.Printf("Container %s is now running", containerName)
				return nil
			}
		} else {
			log.Println("Container not found, waiting...")
		}

		time.Sleep(1 * time.Second)
	}
}

func VerifyIfContainerExists() bool {
	log.Printf("Verifying if container %s exists", c.ComposeContainerName)

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Printf("Failed to create Docker client: %v", err)
		return false
	}
	defer cli.Close()

	_, err = cli.ContainerInspect(context.Background(), c.ComposeContainerName)
	exists := err == nil

	if exists {
		log.Printf("Container %s exists", c.ComposeContainerName)
	} else {
		log.Printf("Container %s does not exist: %v", c.ComposeContainerName, err)
	}

	return exists
}
