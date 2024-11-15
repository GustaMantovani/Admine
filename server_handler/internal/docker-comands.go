package internal

import "os/exec"

func StartServerDockerCompose(composeDirectory string) {
	cmd := exec.Command("docker", "compose", "up", "-d")
	cmd.Dir = composeDirectory
}
