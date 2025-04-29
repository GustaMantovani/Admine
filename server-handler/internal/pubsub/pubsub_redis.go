package pubsub

import (
	"context"
	"log"
	"server-handler/internal/handler"

	"github.com/redis/go-redis/v9"
)

type PubSubRedis struct {
	channel string
	client  *redis.Client
}

func New(adress, channel string) PubSubRedis {
	rdb := redis.NewClient(
		&redis.Options{Addr: adress},
	)

	return PubSubRedis{
		channel: channel,
		client:  rdb,
	}
}

func (ps PubSubRedis) SendMessage(message string) {
	ps.client.Publish(context.Background(), ps.channel, message)
}

func (ps PubSubRedis) ListenForMessages() {
	subscriber := ps.client.Subscribe(context.Background(), ps.channel)

	for {
		msg, err := subscriber.ReceiveMessage(context.Background())

		if err != nil {
			log.Fatal("Erro ao receber mensagem do canal redis: ", err)
		}

		handler.ManageCommand(msg.Payload)
	}
}
