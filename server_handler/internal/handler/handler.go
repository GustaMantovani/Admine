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
	"time"
)

var c = config.GetInstance()

func ManageCommand(msg models.Message, ps pubsub.PubSubInterface) error {
	log.Printf("Managing command with tag: %s", msg.Tags[0])
	if msg.Tags[0] == "server_up" {
		serverUp(ps)
	} else if msg.Tags[0] == "server_down" {
		serverDown(ps)
	} else if msg.Tags[0] == "command" {
		command(ps, msg.Msg)
	} else {
		log.Printf("Received invalid tag: %s", msg.Tags[0])
	}

	return nil
}

func serverUp(ps pubsub.PubSubInterface) {

	if docker.VerifyIfContainerExists() && docker.IsContainerRunning() {
		log.Println("Container already exists and is running, no need to start again")
		nodeID, err := docker.GetZeroTierNodeID(c.ComposeContainerName)
		if err != nil {
			log.Printf("Failed to get ZeroTier Node ID: %v", err)
			return
		}
		log.Printf("Server is already up with ZeroTier Node ID: %s", nodeID)
		msg := models.NewMessage(nodeID, []string{"server_up"})
		ps.SendMessage(msg.ToString(), c.SenderChannel)
		return
	}

	if !docker.VerifyIfContainerExists() {
		log.Println("Waiting for container to build and start...")
	}

	err := minecraftserver.StartServerCompose()
	if err != nil {
		log.Printf("Failed to start server: %v", err)
		return
	}

	err = docker.WaitForBuildAndStart()
	if err != nil {
		log.Printf("Failed to wait for container start: %v", err)
		return
	}

	if !docker.VerifyIfContainerExists() {
		log.Println("Container does not exist after start")
		return
	}

	if !docker.IsContainerRunning() {
		log.Println("Container is not running after start")

		last, err := docker.ReadLastContainerLine()
		if err != nil {
			log.Printf("Error reading container log: %v", err)
		}

		log.Printf("Last container log: %s", last)
		return
	}

	time.Sleep(1 * time.Second)

	log.Println("Getting ZeroTier Node ID...")
	nodeID, err := docker.GetZeroTierNodeID(c.ComposeContainerName)
	if err != nil {
		log.Printf("Failed to get ZeroTier Node ID: %v", err)

		lastLine, err := docker.ReadLastContainerLine()
		if err != nil {
			log.Printf("Error reading last container log: %v", err)
		} else if lastLine != "" {
			log.Printf("Last container log: %s", lastLine)

		}

		return
	}

	time.Sleep(100 * time.Millisecond)

	log.Printf("Server is now up with ZeroTier Node ID: %s", nodeID)
	msg := models.NewMessage(nodeID, []string{"server_up"})
	ps.SendMessage(msg.ToString(), c.SenderChannel)
}

func serverDown(ps pubsub.PubSubInterface) {

	if !docker.VerifyIfContainerExists() {
		log.Println("Container does not exist, cannot stop server")
		return
	}

	if !docker.IsContainerRunning() {
		log.Println("Container is not running, cannot stop server")
		return
	}

	log.Println("Sending stop command to minecraft server...")
	err := commandexecuter.WriteToContainer("/stop")
	if err != nil {
		log.Printf("Failed to send stop command: %v", err)
	}

	log.Println("Waiting for server to save all dimensions...")

	finished := false
	for !finished {
		msg, err := docker.ReadLastContainerLine()
		if err != nil {
			log.Printf("Error reading container log: %v", err)

			if docker.IsContainerRunning() {
				log.Println("Container is still running, but error occurred while reading logs")
			} else {
				log.Println("Container is not running, stopping the process")
				return
			}
		}
		if strings.Contains(msg, "All dimensions are saved") {
			log.Println("Server has saved all dimensions, proceeding with shutdown")
			finished = true
		}
	}

	log.Println("Stopping docker compose...")
	minecraftserver.StopServerCompose()

	if docker.IsContainerRunning() {
		log.Printf("Container is still running after stop command")
		return
	}

	msg := models.NewMessage("", []string{"server_down"})
	ps.SendMessage(msg.ToString(), c.SenderChannel)
}

func command(ps pubsub.PubSubInterface, message string) {
	log.Printf("Processing command: %s", message)

	if !docker.VerifyIfContainerExists() {
		log.Println("Container does not exist, cannot execute command")
		return
	}

	if !docker.IsContainerRunning() {
		log.Println("Container is not running, cannot execute command")
		return
	}

	err := commandexecuter.WriteToContainer(message)
	if err != nil {
		log.Printf("Failed to execute command: %v", err)
		return
	}

	nodeID, err := docker.GetZeroTierNodeID(c.ComposeContainerName)
	if err != nil {
		log.Printf("Failed to get ZeroTier Node ID after command execution: %v", err)
		return
	}

	msg := models.NewMessage(nodeID, []string{"command_executed"})
	ps.SendMessage(msg.ToString(), c.SenderChannel)
}
