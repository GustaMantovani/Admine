package pubsub

import (
	"context"

	"github.com/redis/go-redis/v9"
)

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

func createRedisClient(addr string) redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	return *rdb
}

func createPubSub(data RedisPubSubMetadata, ctx context.Context, client redis.Client) redis.PubSub {
	return *client.Subscribe(ctx, data.channel)
}

func CreateSubscriber(channel, addr string) RedisPubSubSubscriber {
	client := createRedisClient(addr)

	return RedisPubSubSubscriber{
		channel: channel,
		client:  client,
		context: context.Background(),
	}
}

func (sub RedisPubSubSubscriber) SendMessage(message string) {
	sub.client.Publish(sub.context, sub.channel, message)
}
