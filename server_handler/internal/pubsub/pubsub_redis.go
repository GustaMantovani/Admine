package pubsub

import (
	"context"
	"encoding/json"
	"log"
	// "server_handler/internal/handler"
	"server_handler/internal/message"

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

// Listen pubsub for messages in format of the struct Message from internal/message
// and send then to a channel from parameter
func (ps PubSubRedis) ListenForMessages(msgChannel chan message.Message) {
	subscriber := ps.client.Subscribe(context.Background(), ps.channel)

	for {
		msg, err := subscriber.ReceiveMessage(context.Background())

		log.Println("Mensagem recebida pelo pubsub redis: ", msg)

		if err != nil {
			log.Fatal("Erro ao receber mensagem do canal redis: ", err)
		}

		var m message.Message

		err = json.Unmarshal([]byte(msg.Payload), &m)

		if err != nil {
			log.Println("erro: ", err)
		}

		msgChannel <- m
	}
}
