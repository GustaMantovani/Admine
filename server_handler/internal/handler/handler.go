package handler

import (
	"log"
	commandhandler "server_handler/internal/command_handler"
	"server_handler/internal/config"
	"server_handler/internal/docker"
	"server_handler/internal/pubsub"
	"server_handler/internal/server"
	"strings"
)

func ManageCommand(tag, message string, ps pubsub.PubSubInterface) error {
	c := config.GetInstance()

	if tag == "start_server" {

		server.StartServerCompose()
		ps.SendMessage("Starting server", c.SenderChannel)
		ps.SendMessage(docker.GetZeroTierNodeID(c.ComposeContainerName), c.SenderChannel)
		log.Println("Start server")

	} else if tag == "stop_server" {

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
		ps.SendMessage("Server stopped", c.SenderChannel)
		ps.SendMessage(docker.GetZeroTierNodeID(c.ComposeContainerName), c.SenderChannel)

		log.Println("Stop server")

	} else if tag == "command" {

		commandhandler.WriteToContainer(message)
		ps.SendMessage(docker.GetZeroTierNodeID(c.ComposeContainerName), c.SenderChannel)

		log.Println("Send a command to the server: ", message)

	} else {

		ps.SendMessage("Invalid tag.", c.SenderChannel)

		log.Println("Received a invalid tag")

	}

	return nil
}
