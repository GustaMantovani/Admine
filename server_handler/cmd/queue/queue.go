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
	log.Println("Starting queue listener on channels:", config.GetInstance().ConsumerChannel)
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

	go ps.ListenForMessages(c.ConsumerChannel, mc)

	for msg := range mc {
		log.Printf("Received message with tags: %v", msg.Tags)
		if len(msg.Tags) > 0 {
			err := handler.ManageCommand(msg, ps)
			if err != nil {
				log.Printf("Error handling command: %v", err)
			}
		} else {
			log.Println("Received message with no tags, ignoring")
		}
	}
}
