package pubsub

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

var serverChannel string

type RedisPubSubMetadata struct {
	Channel string
	Addr    string
}

// Faz a conexão com o pubsub e envia mensagens para um canal associado ao tipo
type RedisPubSubSubscriber struct {
	channel string
	Client  redis.Client
	context context.Context
}

func GetConfigServerChannelFromDotEnv(envVarName string) RedisPubSubMetadata {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return RedisPubSubMetadata{
		Channel: os.Getenv(envVarName),
		Addr:    os.Getenv("REDIS_URL") + ":" + os.Getenv("REDIS_PORT"),
	}
}

func CreateRedisClient(addr string) redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	return *rdb
}

func CreatePubSub(channel string, ctx context.Context, client redis.Client) redis.PubSub {
	return *client.Subscribe(ctx, channel)
}

func CreateSubscriber(addr, channel string) RedisPubSubSubscriber {
	client := CreateRedisClient(addr)

	return RedisPubSubSubscriber{
		channel: channel,
		Client:  client,
		context: context.Background(),
	}
}

func (sub RedisPubSubSubscriber) SendMessage(message string) {
	log.Println("Canal: ", sub.channel)
	sub.Client.Publish(sub.context, sub.channel, message)
}
