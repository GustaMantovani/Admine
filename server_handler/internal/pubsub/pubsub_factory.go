package pubsub

import "server_handler/internal/config"

// Returns a concrete Pubsub by config.Config pubsub field
func PubSubFactoryCreate() PubSubInterface {
	var c = config.GetInstance()
	var address = c.Host + ":" + c.Port

	var pubSubTypes = map[string]PubSubInterface{
		"redis": New(address),
	}

	return pubSubTypes[c.Pubsub]
}
