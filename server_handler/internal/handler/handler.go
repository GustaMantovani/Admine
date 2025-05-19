package handler

import (
	"log"
	"server_handler/internal/commandexecuter"
	"server_handler/internal/config"
	"server_handler/internal/docker"
	"server_handler/internal/pubsub"
	"server_handler/internal/server"
	"strings"
)

var c = config.GetInstance()

func ManageCommand(tag, message string, ps pubsub.PubSubInterface) error {
	if tag == "server_up" {
		serverUp(ps)
	} else if tag == "server_down" {
		serverDown(ps)
	} else if tag == "command" {
		command(ps, message)
	} else {
		ps.SendMessage("Invalid tag.", c.SenderChannel)

		log.Println("Received a invalid tag")
	}

	return nil
}

func serverUp(ps pubsub.PubSubInterface) {
	server.StartServerCompose()
	ps.SendMessage("Starting server", c.SenderChannel)
	ps.SendMessage(docker.GetZeroTierNodeID(c.ComposeContainerName), c.SenderChannel)
	log.Println("Server up")
}

func serverDown(ps pubsub.PubSubInterface) {
	commandexecuter.WriteToContainer("/stop")
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

	log.Println("Stop server")
}

func command(ps pubsub.PubSubInterface, message string) {
	commandexecuter.WriteToContainer(message)
	ps.SendMessage(docker.GetZeroTierNodeID(c.ComposeContainerName), c.SenderChannel)

	log.Println("Send a command to the server: ", message)
}
