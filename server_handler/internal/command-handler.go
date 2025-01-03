package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"server/handler/pubsub"
)

func RunCommandHandler(containerName string) {
	pubsubAddr := pubsub.GetConfigServerChannelFromDotEnv("REDIS_COMMAND_CHANNEL")
	fmt.Println(pubsubAddr.Addr)
	receiver := pubsub.CreatePubSub(pubsubAddr.Channel, context.Background(), pubsub.CreateRedisClient(pubsubAddr.Addr))
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

		fmt.Println(containerName)

		WriteToContainerByName(containerName, commmand)
	}
}
