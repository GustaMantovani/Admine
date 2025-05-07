package handler

import (
	"errors"
	"log"
	"os/exec"
	commandhandler "server_handler/internal/command_handler"
	"server_handler/internal/config"
	"server_handler/internal/docker"
	"server_handler/internal/pubsub"
	"server_handler/internal/server"
	"strings"
)

func ManageCommand(command string, ps pubsub.PubSubInterface) error {
	c := config.GetInstance()
	if command == "start_server" {
		server.StartServerCompose()
		log.Println("Start server")
		ps.SendMessage("Starting server", c.SenderChannel)
	} else if command == "stop_server" {
		commandhandler.WriteToContainer("/stop")
		ps.SendMessage("Stopping server", c.SenderChannel)
		sair := false
		for sair {
			msg, err := docker.ReadLastContainerLine()
			if err != nil {
				log.Println("Erro ao ler a última linha do container do servidor: ", err.Error())
			}
			if strings.Contains(msg, "All dimensions are saved") {
				sair = true
			}
		}
		server.StopServerCompose()
		log.Println("Stop server")
		ps.SendMessage("Server stopped", c.SenderChannel)
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
