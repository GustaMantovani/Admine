package pubsub

import (
	"context"
	"encoding/json"
	"log"
	"server_handler/internal/message"

	"github.com/redis/go-redis/v9"
)

type PubSubRedis struct {
	client *redis.Client
}

func New(adress string) PubSubRedis {
	rdb := redis.NewClient(
		&redis.Options{Addr: adress},
	)

	return PubSubRedis{
		client: rdb,
	}
}

func (ps PubSubRedis) SendMessage(message, channel string) {
	ps.client.Publish(context.Background(), channel, message)
}

// Listen pubsub for messages in format of the struct Message from internal/message
// and send then to a channel from parameter
func (ps PubSubRedis) ListenForMessages(channel string, msgChannel chan message.Message) {
	subscriber := ps.client.Subscribe(context.Background(), channel)
	ch := subscriber.Channel()

	for msg := range ch {
		var m message.Message

		err := json.Unmarshal([]byte(msg.Payload), &m)

		if err != nil {
			log.Println("erro: ", err)
		}

		msgChannel <- m
	}
}
