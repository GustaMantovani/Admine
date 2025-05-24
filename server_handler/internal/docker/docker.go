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

	// Search container by name
	containerID := getContainerId(containerName, cli)
	if containerID == "" {
		log.Printf("Container %s not found", containerName)
		return "", fmt.Errorf("container not found: %s", containerName)
	}

	// Get logs
	// log.Printf("Fetching logs from container: %s", containerID)
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

func GetZeroTierNodeID(containerName string) (string, error) {
	commandArgs := []string{"docker", "exec", "-i", containerName, "/bin/bash", "-c", "zerotier-cli info"}
	log.Printf("[ZeroTier] Executing: %v", commandArgs)

	cmd := exec.Command(commandArgs[0], commandArgs[1:]...)
	output, err := cmd.CombinedOutput()
	outputStr := string(output)

	if err != nil {
		log.Printf("[ZeroTier] Command failed: %v | Output: %s", err, outputStr)
		return "", fmt.Errorf("failed to execute zerotier-cli info: %w", err)
	}

	parts := strings.Split(outputStr, " ")
	if len(parts) < 3 {
		log.Printf("[ZeroTier] Unexpected output: %s", outputStr)
		return "", fmt.Errorf("invalid ZeroTier CLI output format: expected at least 3 parts, got %d", len(parts))
	}

	nodeID := strings.TrimSpace(parts[2])
	return nodeID, nil
}

func WaitForBuildAndStart() error {
	log.Println("Waiting for container to build and start...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Create Docker client
	cli, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		log.Printf("Error creating Docker client: %v", err)
		return fmt.Errorf("error creating Docker client: %w", err)
	}

	containerName := c.ComposeContainerName
	log.Printf("Waiting for container: %s", containerName)

	err = waitForContainerStart(cli, containerName)
	if err != nil {
		log.Printf("Container start wait failed: %v", err)
		return fmt.Errorf("container start wait failed: %w", err)
	}

	// Verify if container exists
	_, err = cli.ContainerInspect(ctx, containerName)
	if err != nil {
		log.Printf("Container not found after start: %v", err)
		return fmt.Errorf("container not found after start: %w", err)
	}

	// log.Printf("Container %s is up and running", containerName)
	log.Println("Container is up and running")
	return nil
}

func getContainerId(containerName string, cli *client.Client) string {
	// log.Printf("Searching for container with name: %s", containerName)

	containers, err := cli.ContainerList(context.Background(), container.ListOptions{All: true})
	if err != nil {
		log.Printf("Failed to list containers: %v", err)
		return ""
	}

	var containerID string
	for _, c := range containers {
		if slices.Contains(c.Names, "/"+containerName) {
			containerID = c.ID
			// log.Printf("Found container %s with ID: %s", containerName, containerID)
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
			// log.Printf("Container status: %s", container.State)
			if container.State == "running" {
				log.Printf("Container %s is now running", containerName)
				return nil
			}
		} else {
			// log.Println("Container not found, waiting...")
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

func IsContainerRunning() bool {
	log.Printf("Checking if container %s is running", c.ComposeContainerName)

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Printf("Failed to create Docker client: %v", err)
		return false
	}
	defer cli.Close()

	containerInfo, err := cli.ContainerInspect(context.Background(), c.ComposeContainerName)
	if err != nil {
		log.Printf("Container %s not found: %v", c.ComposeContainerName, err)
		return false
	}

	isRunning := containerInfo.State.Running
	if isRunning {
		log.Printf("Container %s is running", c.ComposeContainerName)
	} else {
		log.Printf("Container %s exists but is not running (status: %s)", c.ComposeContainerName, containerInfo.State.Status)
	}

	return isRunning
}
