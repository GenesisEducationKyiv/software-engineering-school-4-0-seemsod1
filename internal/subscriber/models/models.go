package subscribermodels

import (
	"encoding/json"
	"fmt"
	"time"
)

// Subscriber is a struct that represents a subscriber to the newsletter
type Subscriber struct {
	Email     string `gorm:"unique;not null"  json:"email"`
	Timezone  int    `json:"timezone"`
	CreatedAt time.Time
}

// CommandData is a struct that represents the data sent as a command to the subscriber service
type CommandData struct {
	Command string  `json:"command"`
	Payload Payload `json:"payload"`
}

// Payload is a struct that represents the payload of a command
type Payload struct {
	Email    string `json:"email"`
	Timezone int    `json:"timezone"`
}

// ReplyData is a struct that represents the data sent as a reply to a command
type ReplyData struct {
	Command string  `json:"command"`
	Reply   string  `json:"reply"`
	Payload Payload `json:"payload"`
}

func SerializeReplyData(data ReplyData) (string, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to serialize message: %w", err)
	}
	return string(jsonData), nil
}
