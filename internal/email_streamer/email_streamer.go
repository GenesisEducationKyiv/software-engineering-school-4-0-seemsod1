package emailstreamer

import (
	"context"
	"encoding/binary"
	"log"
	"time"

	"github.com/seemsod1/api-project/internal/notifier"

	"github.com/seemsod1/api-project/internal/logger"
	"github.com/segmentio/kafka-go"
)

const (
	RecoverTime = 1 * time.Minute
	PeriodTime  = 1 * time.Minute
	BatchSize   = 100
)

type EmailStreamer struct {
	EventStorage    EventRepo
	StreamerStorage StreamerRepo
	KafkaProducer   *kafka.Writer
	Logger          *logger.Logger
}

type (
	EventRepo interface {
		GetEvents(offset uint, limit int) ([]notifier.Event, error)
	}
	StreamerRepo interface {
		ChangeOffset(msg Streamer) error
		GetOffset(topic string, partition int) (uint, error)
	}
)

func NewEmailStreamer(eventStorage EventRepo,
	streamerStorage StreamerRepo,
	kafkaProducer *kafka.Writer,
	logg *logger.Logger,
) *EmailStreamer {
	return &EmailStreamer{
		EventStorage:    eventStorage,
		StreamerStorage: streamerStorage,
		KafkaProducer:   kafkaProducer,
		Logger:          logg,
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
	off, err := es.StreamerStorage.GetOffset(es.KafkaProducer.Topic, 1)
	if err != nil {
		es.Logger.Errorf("Error retrieving last offset: %v", err)
		time.Sleep(RecoverTime)
		return
	}
	events, err := es.EventStorage.GetEvents(off, BatchSize)
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
		if err = es.StreamerStorage.ChangeOffset(streamMsg); err != nil {
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
