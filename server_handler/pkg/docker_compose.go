package pkg

import (
	"os"
	"os/exec"
)

// DockerCompose is a wrapper for docker compose commands
type DockerCompose struct {
	File string
}

// NewDockerCompose creates a Compose instance using global logger
func NewDockerCompose(file string) *DockerCompose {
	return &DockerCompose{
		File: file,
	}
}

func (dc *DockerCompose) run(args ...string) error {
	baseArgs := []string{"compose"}
	if dc.File != "" {
		baseArgs = append(baseArgs, "-f", dc.File)
	}

	cmdArgs := append(baseArgs, args...)

	Logger.Info("Running command: docker %v", cmdArgs)

	cmd := exec.Command("docker", cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		Logger.Error("Command failed: docker %v | %v", cmdArgs, err)
		return err
	}

	Logger.Info("Command succeeded: docker %v", cmdArgs)

	return nil
}

// Up starts services if specified; otherwise all
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

// Down stops services if specified; otherwise all
func (dc *DockerCompose) Down(services ...string) error {
	args := []string{"down"}
	if len(services) > 0 {
		args = append(args, services...)
	}
	return dc.run(args...)
}

// Stops services if specified; otherwise all
func (dc *DockerCompose) Stop(services ...string) error {
	args := []string{"stop"}
	if len(services) > 0 {
		args = append(args, services...)
	}
	return dc.run(args...)
}

// Ps lists services if specified; otherwise all
func (dc *DockerCompose) Ps(services ...string) error {
	args := []string{"ps"}
	if len(services) > 0 {
		args = append(args, services...)
	}
	return dc.run(args...)
}

// Logs shows logs for services if specified; otherwise all
func (dc *DockerCompose) Logs(services ...string) error {
	args := []string{"logs"}
	if len(services) > 0 {
		args = append(args, services...)
	}
	return dc.run(args...)
}

// Exec runs a command for each specified service
func (dc *DockerCompose) Exec(command []string, services ...string) error {
	if len(services) == 0 {
		Logger.Error("No services specified for Exec")
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

// ExecStructured runs a command for each specified service and returns a map with outputs
func (dc *DockerCompose) ExecStructured(command []string, services ...string) (map[string]string, error) {
	if len(services) == 0 {
		Logger.Info("No services specified for Exec")
		return nil, nil
	}

	results := make(map[string]string)
	for _, service := range services {
		args := append([]string{"exec", service}, command...)
		cmd := exec.Command("docker", append([]string{"compose", "-f", dc.File}, args...)...)

		output, err := cmd.CombinedOutput()
		results[service] = string(output) // store output

		if err != nil {
			return results, err // return partial results if there was an error
		}
	}

	return results, nil
}
