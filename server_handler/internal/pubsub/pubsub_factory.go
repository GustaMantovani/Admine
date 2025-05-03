package pubsub

var pubSubTypes = map[string]PubSubInterface{
	"redis":        New("127.0.0.1:6379", "teste"),
	"redis-sender": New("127.0.0.1:6379", "teste-sender"),
}

func PubSubFactoryCreate(pubsub string) PubSubInterface {
	return pubSubTypes[pubsub]
}
