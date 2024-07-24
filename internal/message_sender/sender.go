package messagesender

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/VictoriaMetrics/metrics"
	"net/smtp"
	"time"

	"go.uber.org/zap"

	"github.com/seemsod1/api-project/pkg/logger"

	"github.com/segmentio/kafka-go"

	"github.com/jordan-wright/email"
)

const (
	numWorkers  = 4
	serviceName = "message_sender"
)

var (
	emailSentSuccessfullyTotal = metrics.NewCounter("email_sent_successfully_total")
	emailSentWithErrorsTotal   = metrics.NewCounter("email_sent_with_errors_total")
	emailSendDurationSummary   = metrics.NewSummary("email_sent_duration_seconds")
	messagesConsumedTotal      = metrics.NewCounter("messages_consumed_total")
)

type SMTPEmailSender struct {
	From         string
	Pool         *email.Pool
	EventStorage eventStorage
	KafkaReader  *kafka.Reader
	Logger       *logger.Logger
}

type eventStorage interface {
	ConsumeEvent(event EventProcessed) error
	CheckEventProcessed(id int) (bool, error)
}

func NewSMTPEmailSender(
	cfg EmailSenderConfig,
	kafkaReader *kafka.Reader,
	storage eventStorage,
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

		var traceID string
		for _, h := range m.Headers {
			if h.Key == "trace_id" {
				traceID = string(h.Value)
				break
			}
		}

		ctx = context.WithValue(ctx, logger.TraceIDKey, traceID)
		ctx = context.WithValue(ctx, logger.ServiceNameKey, serviceName)

		if err != nil {
			s.Logger.WithContext(ctx).Error("reading message", zap.Error(err))
			continue
		}

		data, err := deserializeData(m.Value)
		if err != nil {
			s.Logger.WithContext(ctx).Error("deserializing message", zap.Error(err))
			continue
		}
		e := email.NewEmail()
		e.From = s.From
		e.To = []string{data.Recipient}
		e.Subject = "Currency rate notification: USD to UAH"
		e.Text = []byte(data.Message)

		s.Logger.WithContext(ctx).Info("Sending message to", zap.String("email", data.Recipient))

		isProcessed, err := s.EventStorage.CheckEventProcessed(int(binary.BigEndian.Uint64(m.Key)))
		if err != nil {
			s.Logger.WithContext(ctx).Error("checking if event is processed", zap.Error(err))
			continue
		}
		if isProcessed {
			s.Logger.WithContext(ctx).Warn("Event is already processed")
			if err = s.KafkaReader.CommitMessages(ctx, m); err != nil {
				s.Logger.WithContext(ctx).Error("committing message", zap.Error(err))
				continue
			}
		}
		startTime := time.Now()
		if err = s.Send(e); err != nil {
			s.Logger.WithContext(ctx).Warn("sending email", zap.Error(err))
			emailSentWithErrorsTotal.Inc()
			continue
		}
		elapsedTime := time.Since(startTime).Seconds()
		emailSendDurationSummary.Update(elapsedTime)

		emailSentSuccessfullyTotal.Inc()

		event := EventProcessed{
			ID:   uint(binary.BigEndian.Uint64(m.Key)),
			Data: string(m.Value),
		}
		if err = s.EventStorage.ConsumeEvent(event); err != nil {
			s.Logger.WithContext(ctx).Error("consuming event", zap.Error(err))
		}

		if err = s.KafkaReader.CommitMessages(ctx, m); err != nil {
			s.Logger.WithContext(ctx).Error("committing message", zap.Error(err))
		}
		messagesConsumedTotal.Inc()
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
