package docker

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strings"
)

// DockerCompose is a wrapper for docker compose commands
type DockerCompose struct {
	File string
}

// NewDockerCompose creates a DockerCompose instance
func NewDockerCompose(file string) *DockerCompose {
	return &DockerCompose{File: file}
}

func (dc *DockerCompose) run(args ...string) error {
	baseArgs := []string{"compose"}
	if dc.File != "" {
		baseArgs = append(baseArgs, "-f", dc.File)
	}

	cmdArgs := append(baseArgs, args...)

	slog.Info("Running command", "cmd", "docker", "args", cmdArgs)

	cmd := exec.Command("docker", cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		slog.Error("Command failed", "cmd", "docker", "args", cmdArgs, "error", err)
		return err
	}

	slog.Info("Command succeeded", "cmd", "docker", "args", cmdArgs)

	return nil
}

// Up starts services; pass detach=true for -d flag
func (dc *DockerCompose) Up(detach bool, services ...string) error {
	args := []string{"up"}
	if detach {
		args = append(args, "-d")
	}
	if len(services) > 0 {
		args = append(args, services...)
	}
	return dc.run(args...)
}

// Down stops and removes services
func (dc *DockerCompose) Down(services ...string) error {
	args := []string{"down"}
	if len(services) > 0 {
		args = append(args, services...)
	}
	return dc.run(args...)
}

// Stop stops services without removing them
func (dc *DockerCompose) Stop(services ...string) error {
	args := []string{"stop"}
	if len(services) > 0 {
		args = append(args, services...)
	}
	return dc.run(args...)
}

// Ps lists services
func (dc *DockerCompose) Ps(services ...string) error {
	args := []string{"ps"}
	if len(services) > 0 {
		args = append(args, services...)
	}
	return dc.run(args...)
}

// Logs shows logs for services
func (dc *DockerCompose) Logs(services ...string) error {
	args := []string{"logs"}
	if len(services) > 0 {
		args = append(args, services...)
	}
	return dc.run(args...)
}

// ReadLastServiceLogs returns the last n log lines for the given services (all if none provided)
func (dc *DockerCompose) ReadLastServiceLogs(n uint, services ...string) ([]string, error) {
	baseArgs := []string{"compose"}
	if dc.File != "" {
		baseArgs = append(baseArgs, "-f", dc.File)
	}

	cmdArgs := append(baseArgs, "logs", "--tail", fmt.Sprintf("%d", n))
	if len(services) > 0 {
		cmdArgs = append(cmdArgs, services...)
	}

	slog.Info("Running command", "cmd", "docker", "args", cmdArgs)

	cmd := exec.Command("docker", cmdArgs...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		rawOutput := strings.TrimSpace(string(output))
		slog.Error("Command failed", "cmd", "docker", "args", cmdArgs, "error", err, "output", rawOutput)
		if rawOutput != "" {
			return nil, fmt.Errorf("%w: %s", err, rawOutput)
		}
		return nil, err
	}

	rawOutput := strings.TrimSpace(string(output))
	if rawOutput == "" {
		return []string{}, nil
	}

	lines := strings.Split(rawOutput, "\n")
	slog.Info("Command succeeded", "cmd", "docker", "args", cmdArgs)

	return lines, nil
}

// Exec runs a command for each specified service
func (dc *DockerCompose) Exec(command []string, services ...string) error {
	if len(services) == 0 {
		slog.Error("No services specified for Exec")
		return nil
	}

	for _, service := range services {
		args := append([]string{"exec", service}, command...)
		if err := dc.run(args...); err != nil {
			return err
		}
	}
	return nil
}

// ExecStructured runs a command for each specified service and returns a map[service]output
func (dc *DockerCompose) ExecStructured(command []string, services ...string) (map[string]string, error) {
	if len(services) == 0 {
		slog.Info("No services specified for Exec")
		return nil, nil
	}

	results := make(map[string]string)
	for _, service := range services {
		args := append([]string{"exec", service}, command...)
		cmd := exec.Command("docker", append([]string{"compose", "-f", dc.File}, args...)...)

		output, err := cmd.CombinedOutput()
		results[service] = string(output)

		if err != nil {
			return results, err
		}
	}

	return results, nil
}
