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

func ManageCommand(msg models.Message, ps pubsub.PubSubInterface) error {
	var c = config.GetInstance()
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
	var c = config.GetInstance()
	err := minecraftserver.StartServerCompose()
	if err != nil {
		config.GetLogger().Error("error starting server compose: " + err.Error())
		return
	}

	ps.SendMessage("Starting server", c.SenderChannel)
	err = docker.WaitForBuildAndStart()
	if err != nil {
		config.GetLogger().Error("Error during build and start of the container: " + err.Error())

	}

	zeroTierID, err := docker.GetZeroTierNodeID(c.ComposeContainerName)
	if err != nil {
		config.GetLogger().Error("Error getting zerotier ID: " + err.Error())
		ps.SendMessage("Error in server's zerotier", c.SenderChannel)
		return
	}

	msg := models.NewMessage(zeroTierID, []string{"server_up"})
	ps.SendMessage(msg.ToString(), c.SenderChannel)
}

func serverDown(ps pubsub.PubSubInterface) {
	var c = config.GetInstance()
	zerotierId, err := docker.GetZeroTierNodeID(c.ComposeContainerName)

	if err != nil {
		config.GetLogger().Error("Error getting zerotier ID: " + err.Error())
		ps.SendMessage("Error in server's zerotier", c.SenderChannel)
		return
	}

	commandexecuter.WriteToContainer("/stop")
	ps.SendMessage("Stopping server", c.SenderChannel)

	sair := false
	for sair {
		msg, err := docker.ReadLastContainerLine()
		if err != nil {
			config.GetLogger().Error("Error reading last container line: " + err.Error())
			ps.SendMessage("Error reading container logs", c.SenderChannel)
			return
		}
		if strings.Contains(msg, "All dimensions are saved") {
			sair = true
		}
	}

	err = minecraftserver.StopServerCompose()
	if err != nil {
		config.GetLogger().Warn("Error stopping server: " + err.Error())
		msg := models.NewMessage("Error stopping server", []string{"server_down"})
		ps.SendMessage(msg.ToString(), c.SenderChannel)
		return
	}

	msg := models.NewMessage(zerotierId, []string{"server_down"})

	ps.SendMessage(msg.ToString(), c.SenderChannel)

	config.GetLogger().Info("Server stopped")
}

func command(ps pubsub.PubSubInterface, message string) {
	var c = config.GetInstance()
	commandexecuter.WriteToContainer(message)

	zeroTierID, err := docker.GetZeroTierNodeID(c.ComposeContainerName)
	if err != nil {
		config.GetLogger().Error("Error getting zerotier ID: " + err.Error())
		return
	}

	msg := models.NewMessage(zeroTierID, []string{"commands"})
	ps.SendMessage(msg.ToString(), c.SenderChannel)

	config.GetLogger().Info("Sent a command to the server: " + message)
}
