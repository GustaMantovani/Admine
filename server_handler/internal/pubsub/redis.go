package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/GustaMantovani/Admine/server_handler/internal/config"
	"github.com/redis/go-redis/v9"
)

type redisPubSub struct {
	client *redis.Client
	ctx    context.Context
	cancel context.CancelFunc
}

func newRedisPubSub(c config.RedisConfig, ctx context.Context) *redisPubSub {
	ctx, cancel := context.WithCancel(ctx)

	client := redis.NewClient(&redis.Options{
		Addr:     c.Addr,
		Password: c.Password,
		DB:       c.Db,
	})

	return &redisPubSub{
		client: client,
		ctx:    ctx,
		cancel: cancel,
	}
}

func (r *redisPubSub) Publish(topic string, msg *AdmineMessage) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	slog.Debug("Publishing pub/sub message", "topic", topic, "payload", string(data))
	return r.client.Publish(r.ctx, topic, data).Err()
}

func (r *redisPubSub) Subscribe(topics ...string) (<-chan *AdmineMessage, error) {
	ch := make(chan *AdmineMessage)

	sub := r.client.Subscribe(r.ctx, topics...)
	_, err := sub.Receive(r.ctx)
	if err != nil {
		close(ch)
		return nil, fmt.Errorf("failed to subscribe: %w", err)
	}

	go func() {
		defer close(ch)
		for {
			select {
			case <-r.ctx.Done():
				return
			default:
				msg, err := sub.ReceiveMessage(r.ctx)
				if err != nil {
					continue
				}
				slog.Debug("Received pub/sub message", "topic", msg.Channel, "payload", msg.Payload)
				var admMsg AdmineMessage
				if err := json.Unmarshal([]byte(msg.Payload), &admMsg); err != nil {
					continue
				}
				ch <- &admMsg
			}
		}
	}()

	return ch, nil
}

func (r *redisPubSub) Close() error {
	r.cancel()
	return r.client.Close()
}
