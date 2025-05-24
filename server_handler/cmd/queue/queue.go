package queue

import (
	"log"
	"server_handler/internal/config"
	"server_handler/internal/handler"
	"server_handler/internal/models"
	"server_handler/internal/pubsub"
)

/*
RunListenQueue starts listening to pubsub for commands
*/
func RunListenQueue() {
	config := config.GetInstance()
	log.Printf("Starting queue listener. Consumer channels: %v", config.ConsumerChannel)
	listenCommands()
}

/*
listenCommands defines two goroutines:
One to listen to the pubsub and another to send commands to the handler.
*/
func listenCommands() {
	c := config.GetInstance()
	ps := pubsub.PubSubFactoryCreate()

	mc := make(chan models.Message)
	log.Println("Starting pubsub listener and command processor...")

	go ps.ListenForMessages(c.ConsumerChannel, mc)

	for msg := range mc {
		log.Printf("Received message: %+v", msg)
		if len(msg.Tags) > 0 {
			log.Printf("Processing message with tag: %s", msg.Tags[0])
			err := handler.ManageCommand(msg, ps)
			if err != nil {
				log.Printf("Error processing command: %v", err)
			}
		} else {
			log.Println("Received message with no tags, ignoring")
		}
	}
}
