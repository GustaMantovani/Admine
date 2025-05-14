package pubsub

import "server_handler/internal/config"

var c = config.GetInstance()

var pubSubTypes = map[string]PubSubInterface{
	"redis": New("127.0.0.1:6379"),
}

// Returns a concrete Pubsub by config.Config pubsub field
func PubSubFactoryCreate() PubSubInterface {
	return pubSubTypes[c.Pubsub]
}
