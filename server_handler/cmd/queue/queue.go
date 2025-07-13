package queue

import (
	"log"
	"server_handler/internal/config"
	"server_handler/internal/handler"
	"server_handler/internal/models"
	"server_handler/internal/pubsub"
	"strings"
)

/*
Start to listen a pubsub for commands
*/
func RunListenQueue() {
	config.GetLogger().Info("Running queue. Consumer channel: [" + strings.Join(config.GetInstance().ConsumerChannel, " ") + "]")
	listenCommands()
}

/*
Define two threads.

One for listen the pubsub and other to send commands to handler.
*/
func listenCommands() {
	c := config.GetInstance()
	ps := pubsub.PubSubFactoryCreate()

	mc := make(chan models.Message)

	go ps.ListenForMessages(c.ConsumerChannel, mc)

	for msg := range mc {
		log.Println(msg)
		if len(msg.Tags) > 0 {
			log.Println(msg.Tags[0])
			handler.ManageCommand(msg, ps)
		}
	}
}
