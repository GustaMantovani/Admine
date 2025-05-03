package pubsub

import "server_handler/internal/message"

type PubSubInterface interface {
	ListenForMessages(chan message.Message)
	SendMessage(string)
}
