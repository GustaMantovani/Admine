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
	channel string
	addr    string
}

// Faz a conexão com o pubsub e envia mensagens para um canal associado ao tipo
type RedisPubSubSubscriber struct {
	channel string
	client  redis.Client
	context context.Context
}

func getChannelFromDotEnv() string {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv("REDIS_SERVER_CHANNEL")
}

func createRedisClient(addr string) redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	return *rdb
}

func createPubSub(data RedisPubSubMetadata, ctx context.Context, client redis.Client) redis.PubSub {
	return *client.Subscribe(ctx, data.channel)
}

func CreateSubscriber(addr string) RedisPubSubSubscriber {
	client := createRedisClient(addr)

	return RedisPubSubSubscriber{
		channel: getChannelFromDotEnv(),
		client:  client,
		context: context.Background(),
	}
}

func (sub RedisPubSubSubscriber) SendMessage(message string) {
	log.Println("Canal: ", sub.channel)
	sub.client.Publish(sub.context, sub.channel, message)
}
