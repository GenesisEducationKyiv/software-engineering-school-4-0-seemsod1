package models

import "time"

type Subscriber struct {
	Email     string `gorm:"unique;not null"  json:"email"`
	CreatedAt time.Time
}
