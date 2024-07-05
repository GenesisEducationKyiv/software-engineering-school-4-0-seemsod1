package emailstreamer

import (
	"time"

	"github.com/segmentio/kafka-go"
)

func NewKafkaConsumer(kafkaURL, topic, groupID string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{kafkaURL},
		GroupID:        groupID,
		Topic:          topic,
		MinBytes:       10e3,
		MaxBytes:       10e6,
		RetentionTime:  24 * 7 * time.Hour,
		CommitInterval: 0,
	})
}
