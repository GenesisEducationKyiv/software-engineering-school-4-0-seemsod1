package models

import "time"

// Subscriber is a struct that represents a subscriber to the newsletter
type Subscriber struct {
	Email     string `gorm:"unique;not null"  json:"email"`
	CreatedAt time.Time
}
