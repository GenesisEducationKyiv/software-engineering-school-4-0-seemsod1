package broker

import (
	"fmt"
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

func NewKafkaTopic(brokerAddress, topic string, numPartition int) error {
	conn, err := kafka.Dial("tcp", brokerAddress)
	if err != nil {
		return fmt.Errorf("failed to dial leader: %w", err)
	}
	defer conn.Close()

	partitions, err := conn.ReadPartitions()
	if err != nil {
		return fmt.Errorf("failed to read partitions: %w", err)
	}

	for _, p := range partitions {
		if p.Topic == topic {
			return nil
		}
	}

	err = conn.CreateTopics(kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     numPartition,
		ReplicationFactor: 1,
	})
	if err != nil {
		return fmt.Errorf("failed to create topic: %w", err)
	}
	return nil
}

func NewKafkaProducer(kafkaURL, topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP(kafkaURL),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
}
