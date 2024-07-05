package emailstreamer

import (
	"encoding/json"
	"fmt"
	"time"
)

// Event is a struct that represents a message to be sent
type Event struct {
	ID        uint      `gorm:"primaryKey"`
	Data      string    `gorm:"type:text" json:"data"`
	CreatedAt time.Time `json:"created_at"`
}

type EventProcessed struct {
	ID          uint      `gorm:"primaryKey"`
	Data        string    `gorm:"type:text" json:"data"`
	ProcessedAt time.Time `json:"created_at"`
}

type Data struct {
	Recipient string `json:"recipient"`
	Message   string `json:"message"`
}

type Streamer struct {
	ID         uint `gorm:"primaryKey"`
	Topic      string
	Partition  int
	LastOffset uint
}

// SerializeData converts a Message struct to a JSON string, excluding ID, CreatedAt and SentAt
func SerializeData(data Data) (string, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to serialize message: %w", err)
	}
	return string(jsonData), nil
}

// DeserializeData deserializes JSON string to Data struct
func DeserializeData(jsonData []byte) (Data, error) {
	var data Data
	err := json.Unmarshal(jsonData, &data)
	if err != nil {
		return Data{}, err
	}
	return data, nil
}
