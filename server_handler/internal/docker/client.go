package docker

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os/exec"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

func WriteToContainer(ctx context.Context, containerName string, input string) error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	defer cli.Close()

	containerID, err := getContainerID(containerName, cli, ctx)
	if err != nil {
		return err
	}

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

	input = input + "\n"
	_, err = io.WriteString(hijackedResp.Conn, input)
	if err != nil {
		return err
	}

	if closer, ok := hijackedResp.Conn.(interface{ CloseWrite() error }); ok {
		if err := closer.CloseWrite(); err != nil {
			return err
		}
	}

	return nil
}

func ReadLastContainerNLines(containerName string, n uint, ctx context.Context) ([]string, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	defer cli.Close()

	containerID, err := getContainerID(containerName, cli, ctx)
	if err != nil {
		return nil, err
	}

	out, err := cli.ContainerLogs(ctx, containerID, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: false,
		Timestamps: false,
		Tail:       strconv.FormatUint(uint64(n), 10),
	})
	if err != nil {
		return nil, err
	}
	defer out.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, out)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(&buf)
	var lines []string
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 8 && line[0] < 32 {
			line = line[8:]
		}
		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if len(lines) == 0 || (len(lines) == 1 && lines[0] == "") {
		return nil, fmt.Errorf("no lines found")
	}

	return lines, nil
}

func GetZeroTierNodeID(containerName string) (string, error) {
	cmd := exec.Command("docker", "exec", "-i", containerName, "/bin/bash", "-c", "zerotier-cli info")
	output, err := cmd.CombinedOutput()
	if err != nil {
		slog.Error("Failed to execute zerotier-cli info", "error", err)
		return "", fmt.Errorf("failed to get zerotier info: %w", err)
	}

	outputStr := strings.TrimSpace(string(output))
	slog.Debug("ZeroTier info output", "container", containerName, "output", outputStr)

	parts := strings.Split(outputStr, " ")

	if len(parts) < 3 {
		return "", fmt.Errorf("unexpected zerotier-cli output format: %s", outputStr)
	}

	slog.Debug("ZeroTier node ID extracted", "container", containerName, "id", parts[2])
	return parts[2], nil
}

// GetTailscaleNodeKey returns the Tailscale node's public key in "nodekey:XXXX" format by
// executing "tailscale status --json" inside the container. The Tailscale API accepts this
// format directly as a device identifier, so it can be used as-is in API calls.
//
// Returns an error if the device is not fully registered yet. Registration is confirmed by
// waiting for TailscaleIPs to be assigned — the coordination server only assigns IPs after
// the device has been authenticated and is visible in the API.
func GetTailscaleNodeKey(containerName string) (string, error) {
	cmd := exec.Command("docker", "exec", "-i", containerName, "tailscale", "status", "--json")
	output, err := cmd.CombinedOutput()
	if err != nil {
		slog.Error("Failed to execute tailscale status", "error", err)
		return "", fmt.Errorf("failed to get tailscale status: %w", err)
	}

	var status struct {
		BackendState string   `json:"BackendState"`
		TailscaleIPs []string `json:"TailscaleIPs"`
		Self         struct {
			DeviceId string `json:"ID"`
		} `json:"Self"`
	}

	slog.Debug("Tailscale status JSON", "container", containerName, "output", strings.TrimSpace(string(output)))

	if err := json.Unmarshal(output, &status); err != nil {
		return "", fmt.Errorf("failed to parse tailscale status JSON: %w", err)
	}

	if status.Self.DeviceId == "" {
		return "", fmt.Errorf("tailscale public key not yet available")
	}

	// Only return once the device has a Tailscale IP — this confirms it has been
	// registered with the coordination server and is visible in the API.
	if len(status.TailscaleIPs) == 0 {
		slog.Debug("Tailscale device not yet registered (no IP assigned)", "container", containerName, "backend_state", status.BackendState)
		return "", fmt.Errorf("tailscale device not yet registered (backend: %s)", status.BackendState)
	}

	slog.Debug("Tailscale node key extracted", "container", containerName, "key", status.Self.DeviceId, "ips", status.TailscaleIPs)
	return status.Self.DeviceId, nil
}

func WaitForContainerStart(containerName string, timeout time.Duration, ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	cli, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return err
	}
	defer cli.Close()

	_, err = cli.ContainerInspect(ctx, containerName)
	if err != nil {
		return fmt.Errorf("container '%s' not found: %w", containerName, err)
	}

	return waitForContainerRunning(cli, containerName, ctx)
}

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

func waitForContainerRunning(cli *client.Client, containerName string, ctx context.Context) error {
	filter := filters.NewArgs(filters.Arg("name", containerName))

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			containers, err := cli.ContainerList(ctx, container.ListOptions{
				Filters: filter,
				All:     true,
			})
			if err != nil {
				slog.Error("Error listing containers", "error", err)
				continue
			}

			if len(containers) > 0 {
				c := containers[0]
				slog.Info("Container status", "container", containerName, "status", c.State)

				if c.State == "running" {
					slog.Info("Container is now running", "container", containerName)
					return nil
				}
			} else {
				slog.Debug("Container not found, waiting", "container", containerName)
			}
		}
	}
}

func StreamContainerLogs(ctx context.Context, containerName string, onLine func(string)) error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	defer cli.Close()

	out, err := cli.ContainerLogs(ctx, containerName, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: false,
		Follow:     true,
		Tail:       "1",
	})
	if err != nil {
		return err
	}
	defer out.Close()

	scanner := bufio.NewScanner(out)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 8 && line[0] < 32 {
			line = line[8:]
		}
		println(line)
		onLine(line)
	}
	return scanner.Err()
}
