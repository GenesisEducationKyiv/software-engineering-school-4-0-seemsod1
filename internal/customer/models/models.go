package customersmodels

import (
	"encoding/json"
	"fmt"
	"time"
)

type Customer struct {
	ID        uint   `gorm:"primaryKey"`
	Email     string `gorm:"unique"`
	Timezone  int    `json:"timezone"`
	Status    string `gorm:"default:'created'"`
	CreatedAt time.Time
}

type CommandData struct {
	Command string  `json:"command"`
	Payload Payload `json:"payload"`
}

type Payload struct {
	Email    string `json:"email"`
	Timezone int    `json:"timezone"`
}

// SerializeCommandData SerializeData converts a Message struct to a JSON string, excluding ID, CreatedAt and SentAt
func SerializeCommandData(data CommandData) (string, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to serialize message: %w", err)
	}
	return string(jsonData), nil
}

type ReplyData struct {
	Command string  `json:"command"`
	Reply   string  `json:"reply"`
	Payload Payload `json:"payload"`
}

func DeserializeReplyData(data string) (ReplyData, error) {
	var replyData ReplyData
	err := json.Unmarshal([]byte(data), &replyData)
	if err != nil {
		return ReplyData{}, fmt.Errorf("failed to deserialize reply data: %w", err)
	}
	return replyData, nil
}
