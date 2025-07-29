package pubsub

import (
	"context"
	"encoding/json"
	"os"
	"server_handler/internal/config"
	"server_handler/internal/models"

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
func (ps PubSubRedis) ListenForMessages(channels []string, msgChannel chan models.Message) {
	subscriber := ps.client.Subscribe(context.Background(), channels...)
	_, err := subscriber.Receive(context.Background())
	if err != nil {
		config.GetLogger().Error("error connecting to pubsub subscriber: " + err.Error())
		config.CloseLogFile()
		os.Exit(1)
	}

	ch := subscriber.Channel()
	config.GetLogger().Info("listening messages from consumer channels")

	for msg := range ch {
		var m models.Message

		err := json.Unmarshal([]byte(msg.Payload), &m)

		if err != nil {
			config.GetLogger().Warn("error in message json unmarshal process: " + err.Error())
		}

		msgChannel <- m
	}
}
