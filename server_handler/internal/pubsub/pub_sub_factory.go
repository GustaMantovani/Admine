package pubsub

import (
	"fmt"

	"admine.com/server_handler/internal"
	pubsub_impls "admine.com/server_handler/internal/pubsub/pubsub_impls"
)

// CreatePubSub returns a concrete PubSubService based on type
func CreatePubSub(pubSubType string) (PubSubService, error) {
	ctx := internal.Get() // singleton AppContext
	if ctx == nil {
		return nil, fmt.Errorf("AppContext not initialized")
	}

	switch pubSubType {
	case "redis":
		// Assumes Redis configuration exists in AppContext.Config
		cfg := ctx.Config
		return pubsub_impls.NewRedisPubSub(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB), nil

	default:
		return nil, fmt.Errorf("unknown pubsub type: %s", pubSubType)
	}
}
