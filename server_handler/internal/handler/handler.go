package handler

import (
	"errors"
	"log"
	"os/exec"
	"server_handler/internal/server"
	"strings"
)

func ManageCommand(command string) error {
	if command == "start" {
		server.StartServerCompose()
		log.Println("Up server compose")
	} else if command == "stop" {
		server.StopServerCompose()
		log.Println("Down server compose")
	} else {
		return errors.New("invalid command")
	}

	return nil
}

func GetZeroTierNodeID(containerName string) string {
	cmd := exec.Command("docker", "exec", "-i", containerName, "/bin/bash", "-c", "zerotier-cli info")
	output, err := cmd.CombinedOutput()

	if err != nil {
		panic(err)
	}

	outputStr := string(output)

	parts := strings.Split(outputStr, " ")

	return parts[2]
}
