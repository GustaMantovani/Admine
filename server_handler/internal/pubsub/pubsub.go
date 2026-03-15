package pubsub

import (
	"context"
	"fmt"

	"github.com/GustaMantovani/Admine/server_handler/internal/config"
)

// PubSubService is the messaging abstraction used by all components
type PubSubService interface {
	Publish(topic string, msg *AdmineMessage) error
	Subscribe(topics ...string) (<-chan *AdmineMessage, error)
	Close() error
}

// NewRedis creates a Redis-backed PubSubService
func NewRedis(c config.PubSubConfig, ctx context.Context) (PubSubService, error) {
	switch c.Type {
	case "redis":
		return newRedisPubSub(c.Redis, ctx), nil
	default:
		return nil, fmt.Errorf("unknown pubsub type: %s", c.Type)
	}
}
