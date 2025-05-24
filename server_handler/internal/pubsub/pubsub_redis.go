package pubsub

import (
	"context"
	"encoding/json"
	"log"
	"server_handler/internal/models"

	"github.com/redis/go-redis/v9"
)

type PubSubRedis struct {
	client *redis.Client
}

func New(address string) PubSubRedis {
	log.Printf("Creating Redis PubSub client for address: %s", address)
	rdb := redis.NewClient(
		&redis.Options{Addr: address},
	)
	return PubSubRedis{
		client: rdb,
	}
}

func (ps PubSubRedis) SendMessage(message, channel string) {
	log.Printf("Sending message to channel %s: %s", channel, message)
	err := ps.client.Publish(context.Background(), channel, message).Err()
	if err != nil {
		log.Printf("Failed to publish message to channel %s: %v", channel, err)
	} else {
		log.Printf("Message successfully sent to channel %s", channel)
	}
}

// ListenForMessages listens for pubsub messages in the format of the Message struct from internal/models
// and sends them to the provided channel
func (ps PubSubRedis) ListenForMessages(channels []string, msgChannel chan models.Message) {
	log.Printf("Starting to listen for messages on channels: %v", channels)

	subscriber := ps.client.Subscribe(context.Background(), channels...)
	_, err := subscriber.Receive(context.Background())
	if err != nil {
		log.Printf("Error receiving from subscription: %v", err)
	}
	ch := subscriber.Channel()
	log.Printf("Successfully listening on channels: %v", channels)

	for msg := range ch {
		log.Printf("Received message from channel %s: %s", msg.Channel, msg.Payload)
		var m models.Message

		err := json.Unmarshal([]byte(msg.Payload), &m)

		if err != nil {
			log.Printf("Error unmarshaling JSON message: %v", err)
			continue
		}

		log.Printf("Parsed message successfully: %+v", m)
		msgChannel <- m
	}
}
