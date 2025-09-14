package pubsub

import (
	"fmt"

	"admine.com/server_handler/internal/config"
	"admine.com/server_handler/internal/pubsub/pubsub_impls"
)

// CreatePubSub returns a concrete PubSubService based on type
func CreatePubSub(c config.PubSubConfig) (PubSubService, error) {

	switch c.Type {
	case "redis":
		return pubsub.NewRedisPubSub(c.Redis), nil

	default:
		return nil, fmt.Errorf("unknown pubsub type: %s", c.Type)
	}
}
