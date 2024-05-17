package storage

import "gorm.io/gorm"

// DatabaseRepo is an interface that defines the methods that a database repository should implement
type DatabaseRepo interface {
	Connection() *gorm.DB
	AddSubscriber(string) error
	GetSubscribers() ([]string, error)
}
