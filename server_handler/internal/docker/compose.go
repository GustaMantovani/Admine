package docker

import (
	"os/exec"
	"server_handler/internal/config"
	"strings"
)

func ComposeUp() {
	path := config.GetInstance().ComposeAbsPath
	cmd := exec.Command("docker", "compose", "up")
	s := strings.Split(path, "/")
	l := len(s)
	dir := strings.Join(s[1:l-1], "")

	cmd.Path = dir
}
