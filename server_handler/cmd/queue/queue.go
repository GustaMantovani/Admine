package queue

import (
	"log"
	"server_handler/internal/config"
	"server_handler/internal/handler"
	"server_handler/internal/message"
	"server_handler/internal/pubsub"
)

/*
Start to listen a pubsub for commands
*/
func RunListenQueue() {
	log.Println("Running queue. Consumer channel: ", config.GetInstance().ConsumerChannel)
	listenCommands()
}

/*
Define two threads.

One for listen the pubsub and other to send commands to handler.
*/
func listenCommands() {
	psType := "redis"
	ps := pubsub.PubSubFactoryCreate(psType)

	if ps == nil {
		log.Fatal("Tipo de PubSub não existe: ", psType)
	}

	// message channel
	mc := make(chan message.Message)

	go ps.ListenForMessages(mc)

	for msg := range mc {
		log.Println(msg)
		if len(msg.Tags) > 0 {
			handler.ManageCommand(msg.Tags[0], ps)
		}
	}
}
