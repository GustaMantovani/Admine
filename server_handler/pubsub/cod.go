package pubsub

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisPubSubMetadata struct {
	channel string
	addr    string
}

func createRedisClient(data RedisPubSubMetadata) redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: data.addr,
	})

	return *rdb
}

func createPubSub(data RedisPubSubMetadata, ctx context.Context, client redis.Client) redis.PubSub {
	return *client.Subscribe(ctx)
}
