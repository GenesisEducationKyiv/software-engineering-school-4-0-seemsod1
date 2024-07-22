package subscriber

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	subscribermodels "github.com/seemsod1/api-project/internal/subscriber/models"
	subscriberrepo "github.com/seemsod1/api-project/internal/subscriber/repository"
	"github.com/seemsod1/api-project/pkg/logger"
	"github.com/segmentio/kafka-go"
)

type Service struct {
	Database Database
	Producer *kafka.Writer
	Consumer *kafka.Reader
	Logger   *logger.Logger
}

type Database interface {
	AddSubscriber(subscriber subscribermodels.Subscriber) error
	RemoveSubscriber(email string) error
	GetSubscribersWithTimezone(timezone int) ([]string, error)
	GetSubscribers() ([]string, error)
}

func NewService(database Database, producer *kafka.Writer, consumer *kafka.Reader, logg *logger.Logger) *Service {
	return &Service{
		Database: database,
		Producer: producer,
		Consumer: consumer,
		Logger:   logg,
	}
}

func (s *Service) StartReceivingMessages(ctx context.Context) {
	for {
		m, err := s.Consumer.FetchMessage(ctx)

		var traceID string
		for _, h := range m.Headers {
			if h.Key == "trace_id" {
				traceID = string(h.Value)
				break
			}
		}

		ctx = context.WithValue(ctx, logger.TraceIDKey, traceID)

		if err != nil {
			s.Logger.WithContext(ctx).Error("failed to read message")
			continue
		}

		if err = s.processMessage(ctx, m); err != nil {
			s.Logger.WithContext(ctx).Error("failed to process message")
		}
	}
}

func (s *Service) processMessage(ctx context.Context, m kafka.Message) error {
	var data subscribermodels.CommandData
	if err := json.Unmarshal(m.Value, &data); err != nil {
		s.Logger.WithContext(ctx).Error("failed to unmarshal data")
		return err
	}

	responseMessage, err := s.handleCommand(ctx, data)
	if err != nil {
		if err = s.Database.RemoveSubscriber(data.Payload.Email); err != nil {
			s.Logger.WithContext(ctx).Error("failed to remove subscriber")
		}
	}

	repl := subscribermodels.ReplyData{
		Command: data.Command,
		Reply:   responseMessage,
		Payload: data.Payload,
	}

	return s.sendReply(ctx, m, repl)
}

func (s *Service) handleCommand(ctx context.Context, data subscribermodels.CommandData) (string, error) {
	var responseMessage string
	var err error

	switch data.Command {
	case "subscribe_by_email":
		err = s.subscribeByEmail(ctx, data.Payload.Email, data.Payload.Timezone)
		if err != nil {
			if errors.Is(err, subscriberrepo.ErrorDuplicateSubscription) {
				responseMessage = "already_exists"
			} else {
				responseMessage = "failed"
			}
		} else {
			responseMessage = "success"
		}
	default:
		responseMessage = "unknown_command"
		err = fmt.Errorf("unknown command: %s", data.Command)
	}

	return responseMessage, err
}

func (s *Service) sendReply(ctx context.Context, m kafka.Message, repl subscribermodels.ReplyData) error {
	serializedData, err := subscribermodels.SerializeReplyData(repl)
	if err != nil {
		s.Logger.WithContext(ctx).Error("failed to serialize data")
		return err
	}

	if err = s.sendResponse(m.Key, []byte(serializedData)); err != nil {
		s.Logger.WithContext(ctx).Error("failed to send response")
		return err
	}

	if err = s.Consumer.CommitMessages(ctx, m); err != nil {
		s.Logger.WithContext(ctx).Error("failed to commit message")
		return err
	}

	return nil
}

func (s *Service) subscribeByEmail(ctx context.Context, email string, timezone int) error {
	s.Logger.WithContext(ctx).Info("subscribing by email")
	return s.Database.AddSubscriber(subscribermodels.Subscriber{
		Email:    email,
		Timezone: timezone,
	})
}

func (s *Service) sendResponse(key, value []byte) error {
	msg := kafka.Message{
		Key:   key,
		Value: value,
	}
	return s.Producer.WriteMessages(context.Background(), msg)
}
