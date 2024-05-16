package storage

import "gorm.io/gorm"

type DatabaseRepo interface {
	Connection() *gorm.DB
	AddSubscriber(string) error
}
