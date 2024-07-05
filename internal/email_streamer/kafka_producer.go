package emailstreamer

import (
	"context"
	"encoding/binary"
	"fmt"
	"log"
	"time"

	"github.com/seemsod1/api-project/internal/logger"
	"github.com/segmentio/kafka-go"
)

const (
	RecoverTime = 1 * time.Minute
	PeriodTime  = 1 * time.Minute
	BatchSize   = 100
)

type EmailStreamer struct {
	Storage       Repo
	KafkaProducer *kafka.Writer
	Logger        *logger.Logger
}

type (
	Repo interface {
		GetOutboxEvents(offset uint, limit int) ([]Event, error)
		ChangeOffset(msg Streamer) error
		GetOffset(topic string, partition int) (uint, error)
	}
)

func NewEmailStreamer(storage Repo, kafkaProducer *kafka.Writer, logg *logger.Logger) *EmailStreamer {
	return &EmailStreamer{
		Storage:       storage,
		KafkaProducer: kafkaProducer,
		Logger:        logg,
	}
}

func NewKafkaTopic(brokerAddress, topic string) error {
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
		NumPartitions:     1,
		ReplicationFactor: 1,
	})
	if err != nil {
		return fmt.Errorf("failed to create topic: %w", err)
	}
	return nil
}

func NewKafkaWriter(kafkaURL, topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP(kafkaURL),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
}

func (es *EmailStreamer) Process(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Shutting down processing...")
			return
		case <-ticker.C:
			es.processEvents()
		}
	}
}

func (es *EmailStreamer) processEvents() {
	off, err := es.Storage.GetOffset(es.KafkaProducer.Topic, 1)
	if err != nil {
		es.Logger.Errorf("Error retrieving last offset: %v", err)
		time.Sleep(RecoverTime)
		return
	}
	events, err := es.Storage.GetOutboxEvents(off, BatchSize)
	if err != nil {
		es.Logger.Errorf("Error retrieving outbox messages: %v", err)
		time.Sleep(RecoverTime)
		return
	}

	for _, msg := range events {
		key := make([]byte, 8)
		binary.BigEndian.PutUint64(key, uint64(msg.ID))
		if err = publishEvent(es.KafkaProducer, key, []byte(msg.Data)); err != nil {
			es.Logger.Errorf("Failed to write message to Kafka: %v", err)
			continue
		}
		streamMsg := Streamer{
			Topic:      es.KafkaProducer.Topic,
			Partition:  1,
			LastOffset: msg.ID,
		}
		if err = es.Storage.ChangeOffset(streamMsg); err != nil {
			es.Logger.Errorf("Failed to add message to stream: %v", err)
			continue
		}

	}

	es.Logger.Info("All outbox messages processed")

	time.Sleep(PeriodTime)
}

func publishEvent(writer *kafka.Writer, key, message []byte) error {
	conn, err := kafka.DialLeader(context.Background(), "tcp", writer.Addr.String(), writer.Topic, 0)
	if err != nil {
		log.Printf("Failed to dial leader: %v", err)
		return err
	}
	defer conn.Close()

	_, err = conn.WriteMessages(
		kafka.Message{
			Key:   key,
			Value: message,
		},
	)
	if err != nil {
		log.Printf("Failed to write to leader: %v", err)
		return err
	}
	return nil
}
