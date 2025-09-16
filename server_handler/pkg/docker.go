package pkg

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os/exec"
	"slices"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

func WriteToContainer(ctx context.Context, containerName string, input string) error {

	// Create Docker client
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	defer cli.Close()

	// Find container ID by name
	containerID, err := getContainerID(containerName, cli, ctx)
	if err != nil {
		return err
	}

	// Attach to container stdin
	hijackedResp, err := cli.ContainerAttach(ctx, containerID, container.AttachOptions{
		Stream: true,
		Stdin:  true,
		Stdout: false,
		Stderr: false,
	})
	if err != nil {
		return err
	}
	defer hijackedResp.Close()

	// Write input to container
	input = input + "\n"
	_, err = io.WriteString(hijackedResp.Conn, input)
	if err != nil {
		return err
	}

	// Optionally close stdin write
	if closer, ok := hijackedResp.Conn.(interface{ CloseWrite() error }); ok {
		if err := closer.CloseWrite(); err != nil {
			return err
		}
	}

	return nil
}

// ReadLastContainerLine reads the last line of logs from a container
func ReadLastContainerLine(containerName string, ctx context.Context) (string, error) {

	// Create Docker client
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return "", err
	}
	defer cli.Close()

	// Find container ID by name
	containerID, err := getContainerID(containerName, cli, ctx)
	if err != nil {
		return "", err
	}

	// Get container logs
	out, err := cli.ContainerLogs(ctx, containerID, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: false,
		Timestamps: false,
		Tail:       "1", // Read only the last line
	})
	if err != nil {
		return "", err
	}
	defer out.Close()

	// Docker logs come with an 8-byte header per stream
	// We need to skip this header to get the actual content
	var buf bytes.Buffer
	_, err = io.Copy(&buf, out)
	if err != nil {
		return "", err
	}

	scanner := bufio.NewScanner(&buf)
	var lastLine string
	for scanner.Scan() {
		line := scanner.Text()
		// Skip Docker header bytes (usually starts with special characters)
		if len(line) > 8 && line[0] < 32 {
			line = line[8:]
		}
		lastLine = line
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	if lastLine == "" {
		return "", fmt.Errorf("no lines found")
	}

	return lastLine, nil
}

// GetZeroTierNodeID gets the ZeroTier node ID from a container
func GetZeroTierNodeID(containerName string) (string, error) {
	cmd := exec.Command("docker", "exec", "-i", containerName, "/bin/bash", "-c", "zerotier-cli info")
	output, err := cmd.CombinedOutput()
	if err != nil {
		slog.Error("Failed to execute zerotier-cli info", "error", err)
		return "", fmt.Errorf("failed to get zerotier info: %w", err)
	}

	outputStr := strings.TrimSpace(string(output))
	parts := strings.Split(outputStr, " ")

	if len(parts) < 3 {
		return "", fmt.Errorf("unexpected zerotier-cli output format: %s", outputStr)
	}

	return parts[2], nil
}

// WaitForContainerStart waits for a container to start and be in running state
func WaitForContainerStart(containerName string, timeout time.Duration, ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Create Docker client
	cli, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return err
	}
	defer cli.Close()

	// Check if container exists
	_, err = cli.ContainerInspect(ctx, containerName)
	if err != nil {
		return fmt.Errorf("container '%s' not found: %w", containerName, err)
	}

	return waitForContainerRunning(cli, containerName, ctx)
}

// getContainerID finds a container ID by name
func getContainerID(containerName string, cli *client.Client, ctx context.Context) (string, error) {
	containers, err := cli.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		return "", err
	}

	for _, c := range containers {
		if slices.Contains(c.Names, "/"+containerName) {
			return c.ID, nil
		}
	}

	return "", fmt.Errorf("container '%s' not found", containerName)
}

// waitForContainerRunning waits for a container to be in running state
func waitForContainerRunning(cli *client.Client, containerName string, ctx context.Context) error {
	filter := filters.NewArgs(filters.Arg("name", containerName))

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			// List containers
			containers, err := cli.ContainerList(ctx, container.ListOptions{
				Filters: filter,
				All:     true,
			})
			if err != nil {
				slog.Error("Error listing containers", "error", err)
				continue
			}

			if len(containers) > 0 {
				container := containers[0]
				slog.Info("Container status", "container", containerName, "status", container.State)

				if container.State == "running" {
					slog.Info("Container is now running", "container", containerName)
					return nil
				}
			} else {
				slog.Debug("Container not found, waiting", "container", containerName)
			}
		}
	}
}
