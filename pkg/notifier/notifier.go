package notifier

import "time"

// Event is a struct that represents a message to be sent
type Event struct {
	ID        uint      `gorm:"primaryKey"`
	Data      string    `gorm:"type:text" json:"data"`
	CreatedAt time.Time `json:"created_at"`
}
