package pubsub

import (
	"github.com/GustaMantovani/Admine/server_handler/internal/pubsub/models"
)

type PubSubService interface {
	Publish(topic string, msg *models.AdmineMessage) error

	Subscribe(topics ...string) (<-chan *models.AdmineMessage, error)

	Close() error
}
