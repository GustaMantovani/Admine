package pubsub

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

var serverChannel string

type PubSubMetadata struct {
	Channel string
	Addr    string
}

// Faz a conexão com o pubsub e envia mensagens para um canal associado ao tipo
type RedisPubSub struct {
	channel string
	Client  redis.Client
	context context.Context
}

func GetConfigServerChannelFromDotEnv(channelVarName string) PubSubMetadata {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return PubSubMetadata{
		Channel: os.Getenv(channelVarName),
		Addr:    os.Getenv("REDIS_URL"),
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

func CreatePubsub(addr, channel string) RedisPubSub {
	client := CreateRedisClient(addr)

	return RedisPubSub{
		channel: channel,
		Client:  client,
		context: context.Background(),
	}
}

func (sub RedisPubSub) SendMessage(message string) {
	// log.Println("Canal: ", sub.channel)
	sub.Client.Publish(sub.context, sub.channel, message)
}

func GetRedisMetadataFromDotEnv(channelVar string) {}

func ListenChannelForMessages(pubsubChannel, address string, goChannel chan string) {
	subscriber := CreatePubSub(pubsubChannel, context.Background(), CreateRedisClient(address))

	for {
		msg, err := subscriber.ReceiveMessage(context.Background())

		if err != nil {
			log.Fatal("Erro ao receber mensagem do canal redis: ", err)
		}

		goChannel <- msg.Payload
	}
}
