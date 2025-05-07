package pubsub

import "server_handler/internal/message"

type PubSubInterface interface {
	ListenForMessages(string, chan message.Message)
	SendMessage(string, string)
}
