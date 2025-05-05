package handler

import (
	"errors"
	"log"
	"os/exec"
	commandhandler "server_handler/internal/command_handler"
	"server_handler/internal/server"
	"strings"
	"time"
)

func ManageCommand(command string) error {
	if command == "start_server" {
		server.StartServerCompose()
		log.Println("Start server")
	} else if command == "stop_server" {
		commandhandler.WriteToContainer("/stop")
		time.Sleep(5 * time.Second)
		server.StopServerCompose()
		log.Println("Stop server")
	} else if command == "ping" {
		commandhandler.WriteToContainer("/say PONG")
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
