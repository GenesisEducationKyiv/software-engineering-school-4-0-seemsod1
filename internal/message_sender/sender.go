package messagesender

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"net/smtp"

	"github.com/seemsod1/api-project/pkg/logger"

	"github.com/segmentio/kafka-go"

	"github.com/jordan-wright/email"
)

const numWorkers = 4

type SMTPEmailSender struct {
	From         string
	Pool         *email.Pool
	EventStorage EventStorage
	KafkaReader  *kafka.Reader
	Logger       *logger.Logger
}

type EventStorage interface {
	ConsumeEvent(event EventProcessed) error
	CheckEventProcessed(id int) (bool, error)
}

func NewSMTPEmailSender(
	cfg EmailSenderConfig,
	kafkaReader *kafka.Reader,
	storage EventStorage,
	logg *logger.Logger,
) (*SMTPEmailSender, error) {
	p, err := email.NewPool(
		fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		numWorkers,
		smtp.PlainAuth("", cfg.From, cfg.Password, cfg.Host),
	)
	if err != nil {
		return nil, fmt.Errorf("creating email pool: %w", err)
	}
	return &SMTPEmailSender{
		Pool:         p,
		EventStorage: storage,
		KafkaReader:  kafkaReader,
		From:         cfg.From,
		Logger:       logg,
	}, nil
}

func (s *SMTPEmailSender) Send(e *email.Email) error {
	return s.Pool.Send(e, -1)
}

func (s *SMTPEmailSender) StartReceivingMessages(ctx context.Context) {
	for {
		m, err := s.KafkaReader.FetchMessage(ctx)
		if err != nil {
			s.Logger.Warn("reading message", zap.Error(err))
			continue
		}
		data, err := deserializeData(m.Value)
		if err != nil {
			s.Logger.Warn("deserializing message", zap.Error(err))
			continue
		}
		e := email.NewEmail()
		e.From = s.From
		e.To = []string{data.Recipient}
		e.Subject = "Currency rate notification: USD to UAH"
		e.Text = []byte(data.Message)

		s.Logger.Info("Sending message to", zap.String("email", data.Recipient))

		isProcessed, err := s.EventStorage.CheckEventProcessed(int(binary.BigEndian.Uint64(m.Key)))
		if err != nil {
			s.Logger.Warn("checking if event is processed", zap.Error(err))
			continue
		}
		if isProcessed {
			s.Logger.Info("Event is already processed")
			if err = s.KafkaReader.CommitMessages(ctx, m); err != nil {
				s.Logger.Warn("committing message", zap.Error(err))
				continue
			}
		}

		if err = s.Send(e); err != nil {
			s.Logger.Warn("sending email", zap.Error(err))
			continue
		}
		event := EventProcessed{
			ID:   uint(binary.BigEndian.Uint64(m.Key)),
			Data: string(m.Value),
		}
		if err = s.EventStorage.ConsumeEvent(event); err != nil {
			s.Logger.Warn("consuming event", zap.Error(err))
		}

		if err = s.KafkaReader.CommitMessages(ctx, m); err != nil {
			s.Logger.Warn("committing message", zap.Error(err))
		}
	}
}

// deserializeData deserializes JSON string to Data struct
func deserializeData(jsonData []byte) (Data, error) {
	var data Data
	err := json.Unmarshal(jsonData, &data)
	if err != nil {
		return Data{}, err
	}
	return data, nil
}
