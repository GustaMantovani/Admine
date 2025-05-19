package pubsub

import (
	"server_handler/internal/models"
)

type PubSubInterface interface {
	ListenForMessages([]string, chan models.Message)
	SendMessage(string, string)
}
