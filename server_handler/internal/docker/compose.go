package docker

import (
	"os/exec"
)

func ComposeUp() {
	cmd := exec.Command("docker", "compose", "up")

	cmd.Path = "/compose/dir"
	cmd.Run()
}
