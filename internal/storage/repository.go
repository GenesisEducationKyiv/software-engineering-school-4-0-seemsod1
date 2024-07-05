package storage

import (
	"github.com/seemsod1/api-project/internal/models"
)

// DatabaseRepo is an interface that defines the methods that a database repository should implement
type DatabaseRepo interface {
	AddSubscriber(subscriber models.Subscriber) error
	GetSubscribersWithTimezone(timezone int) ([]string, error)
	GetSubscribers() ([]string, error)
}
