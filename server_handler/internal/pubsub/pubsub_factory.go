package pubsub

import (
	"context"
	"fmt"

	"github.com/GustaMantovani/Admine/server_handler/internal/config"
	pubsub_impls "github.com/GustaMantovani/Admine/server_handler/internal/pubsub/pubsub_impls"
)

// CreatePubSub returns a concrete PubSubService based on type
func CreatePubSub(c config.PubSubConfig, ctx context.Context) (PubSubService, error) {

	switch c.Type {
	case "redis":
		return pubsub_impls.NewRedisPubSub(c.Redis, ctx), nil

	default:
		return nil, fmt.Errorf("unknown pubsub type: %s", c.Type)
	}
}
