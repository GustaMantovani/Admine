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
	dir := "/" + strings.Join(s[1:l-1], "/") + "/"

	file := config.GetLogFile()

	cmd.Dir = dir
	cmd.Stdout = file
	go func() {
		err := cmd.Run()

		if err != nil {
			// TODO: upgrade error msg
			config.GetLogger().Error("Error")
		}
	}()
}

func ComposeDown() {
	path := config.GetInstance().ComposeAbsPath
	cmd := exec.Command("docker", "compose", "down")
	s := strings.Split(path, "/")
	l := len(s)
	dir := "/" + strings.Join(s[1:l-1], "/") + "/"

	file := config.GetLogFile()

	cmd.Dir = dir
	cmd.Stdout = file
	go func() {
		err := cmd.Run()

		if err != nil {
			// TODO: upgrade error msg
			config.GetLogger().Error("Error")
		}
	}()
}
