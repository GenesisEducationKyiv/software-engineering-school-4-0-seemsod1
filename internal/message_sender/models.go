package messagesender

import "time"

type EventProcessed struct {
	ID          uint      `gorm:"primaryKey"`
	Data        string    `gorm:"type:text" json:"data"`
	ProcessedAt time.Time `json:"created_at"`
}

// Data is a struct that represents the data of a message
type Data struct {
	Recipient string `json:"recipient"`
	Message   string `json:"message"`
}
