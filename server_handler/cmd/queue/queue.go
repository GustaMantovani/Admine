package queue

import (
	"log"
	"server-handler/internal/pubsub"
)

func RunListenQueue() {
	psType := "redis"
	ps := pubsub.PubSubFactoryCreate(psType)

	if ps == nil {
		log.Fatal("Tipo de PubSub não existe: ", psType)
	}

	ps.ListenForMessages()
}
