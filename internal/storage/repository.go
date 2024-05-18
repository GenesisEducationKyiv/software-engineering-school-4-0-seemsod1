package storage

// DatabaseRepo is an interface that defines the methods that a database repository should implement
type DatabaseRepo interface {
	AddSubscriber(string) error
	GetSubscribers() ([]string, error)
}
