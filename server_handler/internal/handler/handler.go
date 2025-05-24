package handler

import (
	"log"
	"server_handler/internal/commandexecuter"
	"server_handler/internal/config"
	"server_handler/internal/docker"
	"server_handler/internal/minecraftserver"
	"server_handler/internal/models"
	"server_handler/internal/pubsub"
	"strings"
)

var c = config.GetInstance()

func ManageCommand(msg models.Message, ps pubsub.PubSubInterface) error {
	if msg.Tags[0] == "server_up" {
		serverUp(ps)
	} else if msg.Tags[0] == "server_down" {
		serverDown(ps)
	} else if msg.Tags[0] == "command" {
		command(ps, msg.Msg)
	} else {
		ps.SendMessage("Invalid tag.", c.SenderChannel)

		log.Println("Received a invalid tag")
	}

	return nil
}

func serverUp(ps pubsub.PubSubInterface) {
	minecraftserver.StartServerCompose()
	ps.SendMessage("Starting server", c.SenderChannel)
	docker.WaitForBuildAndStart()
	msg := models.NewMessage(docker.GetZeroTierNodeID(c.ComposeContainerName), []string{"server_up"})
	ps.SendMessage(msg.ToString(), c.SenderChannel)
	log.Println("Server up")
}

func serverDown(ps pubsub.PubSubInterface) {
	zerotierId := docker.GetZeroTierNodeID(c.ComposeContainerName)

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

	minecraftserver.StopServerCompose()

	msg := models.NewMessage(zerotierId, []string{"server_down"})

	ps.SendMessage(msg.ToString(), c.SenderChannel)

	log.Println("Stop server")
}

func command(ps pubsub.PubSubInterface, message string) {
	if !docker.VerifyIfContainerExists() {
		log.Println("Container server dont exists")
		return
	}

	commandexecuter.WriteToContainer(message)
	msg := models.NewMessage(docker.GetZeroTierNodeID(c.ComposeContainerName), []string{"commands"})
	ps.SendMessage(msg.ToString(), c.SenderChannel)

	log.Println("Send a command to the server: ", message)
}
