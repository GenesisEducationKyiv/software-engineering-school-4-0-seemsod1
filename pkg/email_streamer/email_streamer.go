package emailstreamer

import (
	"context"
	"encoding/binary"
	"go.uber.org/zap"
	"time"

	"github.com/seemsod1/api-project/pkg/notifier"

	"github.com/seemsod1/api-project/pkg/logger"
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
			es.Logger.Info("Shutting down processing...")
			return
		case <-ticker.C:
			es.processEvents()
		}
	}
}

func (es *EmailStreamer) processEvents() {
	off, err := es.StreamerStorage.GetOffset(es.KafkaProducer.Topic, 1)
	if err != nil {
		es.Logger.Error("Error retrieving last offset", zap.Error(err))
		time.Sleep(RecoverTime)
		return
	}
	events, err := es.EventStorage.GetEvents(off, BatchSize)
	if err != nil {
		es.Logger.Error("Error retrieving outbox messages", zap.Error(err))
		time.Sleep(RecoverTime)
		return
	}

	for _, msg := range events {
		key := make([]byte, 8)
		binary.BigEndian.PutUint64(key, uint64(msg.ID))
		if err = es.publishEvent(key, []byte(msg.Data)); err != nil {
			es.Logger.Error("Failed to write message to broker", zap.Error(err), zap.String("key", string(key)))
			continue
		}
		streamMsg := Streamer{
			Topic:      es.KafkaProducer.Topic,
			Partition:  1,
			LastOffset: msg.ID,
		}
		if err = es.StreamerStorage.ChangeOffset(streamMsg); err != nil {
			es.Logger.Error("Failed to add message to stream", zap.Error(err), zap.Any("message", msg))
			continue
		}

	}

	es.Logger.Info("All outbox messages processed")

	time.Sleep(PeriodTime)
}

func (es *EmailStreamer) publishEvent(key, message []byte) error {
	conn, err := kafka.DialLeader(context.Background(), "tcp", es.KafkaProducer.Addr.String(), es.KafkaProducer.Topic, 0)
	if err != nil {
		es.Logger.Error("Failed to dial leader", zap.Error(err))
		return err
	}
	defer func(conn *kafka.Conn) {
		err = conn.Close()
		if err != nil {
			es.Logger.Error("Failed to close connection with kafka", zap.Error(err))
		}
	}(conn)

	_, err = conn.WriteMessages(
		kafka.Message{
			Key:   key,
			Value: message,
		},
	)
	if err != nil {
		es.Logger.Error("Failed to write to leader", zap.Error(err))
		return err
	}
	return nil
}
