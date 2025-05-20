package pubsub

import "server_handler/internal/config"

var c = config.GetInstance()
var address = c.Host + ":" + c.Port

var pubSubTypes = map[string]PubSubInterface{
	"redis": New(address),
}

// Returns a concrete Pubsub by config.Config pubsub field
func PubSubFactoryCreate() PubSubInterface {
	return pubSubTypes[c.Pubsub]
}
