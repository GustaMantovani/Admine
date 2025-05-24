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
	log.Printf("Processing command with tag: %s", msg.Tags[0])

	if msg.Tags[0] == "server_up" {
		serverUp(ps)
	} else if msg.Tags[0] == "server_down" {
		serverDown(ps)
	} else if msg.Tags[0] == "command" {
		command(ps, msg.Msg)
	} else {
		invalidMsg := models.NewMessage("Invalid tag received", []string{"error"})
		ps.SendMessage(invalidMsg.ToString(), c.SenderChannel)
		log.Printf("Received invalid tag: %s", msg.Tags[0])
	}

	return nil
}

func serverUp(ps pubsub.PubSubInterface) {
	log.Println("Starting server up process...")

	err := minecraftserver.StartServerCompose()
	if err != nil {
		log.Printf("Failed to start server: %v", err)
		errorMsg := models.NewMessage("Failed to start server", []string{"error"})
		ps.SendMessage(errorMsg.ToString(), c.SenderChannel)
		return
	}

	startingMsg := models.NewMessage("Server is starting up", []string{"server_status"})
	ps.SendMessage(startingMsg.ToString(), c.SenderChannel)

	log.Println("Waiting for container to build and start...")
	docker.WaitForBuildAndStart()

	nodeID := docker.GetZeroTierNodeID(c.ComposeContainerName)
	log.Printf("Server started successfully. ZeroTier Node ID: %s", nodeID)

	msg := models.NewMessage(nodeID, []string{"server_up"})
	ps.SendMessage(msg.ToString(), c.SenderChannel)
}

func serverDown(ps pubsub.PubSubInterface) {
	log.Println("Starting server shutdown process...")

	zerotierId := docker.GetZeroTierNodeID(c.ComposeContainerName)
	log.Printf("Retrieved ZeroTier Node ID: %s", zerotierId)

	log.Println("Sending stop command to minecraft server...")
	err := commandexecuter.WriteToContainer("/stop")
	if err != nil {
		log.Printf("Failed to send stop command: %v", err)
	}

	stoppingMsg := models.NewMessage("Server is shutting down", []string{"server_status"})
	ps.SendMessage(stoppingMsg.ToString(), c.SenderChannel)

	log.Println("Waiting for server to save all dimensions...")
	finished := false
	for !finished {
		msg, err := docker.ReadLastContainerLine()
		if err != nil {
			log.Printf("Error reading container log: %v", err)
		}
		if strings.Contains(msg, "All dimensions are saved") {
			log.Println("Server has saved all dimensions, proceeding with shutdown")
			finished = true
		}
	}

	log.Println("Stopping docker compose...")
	minecraftserver.StopServerCompose()

	msg := models.NewMessage(zerotierId, []string{"server_down"})
	ps.SendMessage(msg.ToString(), c.SenderChannel)

	log.Println("Server shutdown completed successfully")
}

func command(ps pubsub.PubSubInterface, message string) {
	log.Printf("Processing command: %s", message)

	if !docker.VerifyIfContainerExists() {
		log.Println("Container does not exist, cannot execute command")
		errorMsg := models.NewMessage("Server container not found", []string{"error"})
		ps.SendMessage(errorMsg.ToString(), c.SenderChannel)
		return
	}

	log.Printf("Executing command on container: %s", message)
	err := commandexecuter.WriteToContainer(message)
	if err != nil {
		log.Printf("Failed to execute command: %v", err)
		errorMsg := models.NewMessage("Failed to execute command", []string{"error"})
		ps.SendMessage(errorMsg.ToString(), c.SenderChannel)
		return
	}

	nodeID := docker.GetZeroTierNodeID(c.ComposeContainerName)
	msg := models.NewMessage(nodeID, []string{"command_executed"})
	ps.SendMessage(msg.ToString(), c.SenderChannel)

	log.Printf("Command executed successfully: %s", message)
}
