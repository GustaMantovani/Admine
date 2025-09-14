package pubsub_impls

import (
	"context"
	"encoding/json"
	"fmt"

	"admine.com/server_handler/internal/config"
	"admine.com/server_handler/internal/pubsub/models"
	"github.com/redis/go-redis/v9"
)

type RedisPubSub struct {
	client *redis.Client
	ctx    context.Context
	cancel context.CancelFunc
}

func NewRedisPubSub(c config.RedisConfig) *RedisPubSub {
	ctx, cancel := context.WithCancel(context.Background())
	client := redis.NewClient(&redis.Options{
		Addr:     c.Addr,
		Password: c.Password,
		DB:       c.Db,
	})

	return &RedisPubSub{
		client: client,
		ctx:    ctx,
		cancel: cancel,
	}
}

func (r *RedisPubSub) Publish(topic string, msg *models.AdmineMessage) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	return r.client.Publish(r.ctx, topic, data).Err()
}

func (r *RedisPubSub) Subscribe(topics ...string) (<-chan *models.AdmineMessage, error) {
	ch := make(chan *models.AdmineMessage)

	pubsub := r.client.Subscribe(r.ctx, topics...)
	_, err := pubsub.Receive(r.ctx)
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
				msg, err := pubsub.ReceiveMessage(r.ctx)
				if err != nil {
					continue
				}
				var admMsg models.AdmineMessage
				if err := json.Unmarshal([]byte(msg.Payload), &admMsg); err != nil {
					continue
				}
				ch <- &admMsg
			}
		}
	}()

	return ch, nil
}

func (r *RedisPubSub) Close() error {
	r.cancel()
	return r.client.Close()
}
