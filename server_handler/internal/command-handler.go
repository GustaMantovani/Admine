package internal

import (
	"context"
	"encoding/json"
	"log"
	"server/handler/pubsub"
)

func RunCommandHandler() {
	receiver := pubsub.CreatePubSub("command", context.Background(), pubsub.CreateRedisClient("localhost:6379"))
	for {
		msg, err := receiver.ReceiveMessage(context.Background())

		if err != nil {
			panic(err)
		}

		var m Message

		err = json.Unmarshal([]byte(msg.Payload), &m)

		if err != nil {
			panic(err)
		}

		commmand := m.Msg + "\n"

		WriteToContainerByName("minecraft-server-mine_server-1", commmand)
		log.Println("Commando recebindo: ", commmand)
	}
}
