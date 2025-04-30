package pubsub

type PubSubInterface interface {
	ListenForMessages()
	SendMessage(string)
}
