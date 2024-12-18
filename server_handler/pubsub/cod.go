package pubsub

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisPubSubMetadata struct {
	channel string
	addr    string
}

type RedisPubSubSubscriber struct {
	channel string
	client  redis.Client
	context context.Context
}

func createRedisClient(data RedisPubSubMetadata) redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: data.addr,
	})

	return *rdb
}

func createPubSub(data RedisPubSubMetadata, ctx context.Context, client redis.Client) redis.PubSub {
	return *client.Subscribe(ctx, data.channel)
}

func createSubscriber(data RedisPubSubMetadata) RedisPubSubSubscriber {
	client := createRedisClient(data)

	return RedisPubSubSubscriber{
		channel: data.channel,
		client:  client,
		context: context.Background(),
	}
}

func (sub RedisPubSubSubscriber) SendMessage(message string) {
	sub.client.Publish(sub.context, sub.channel, message)
}
