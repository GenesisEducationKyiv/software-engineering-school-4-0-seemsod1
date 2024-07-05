package notifier

import (
	"context"
	"fmt"
	"time"

	emailStreamer "github.com/seemsod1/api-project/internal/email_streamer"
	"github.com/segmentio/kafka-go"

	"github.com/jordan-wright/email"
	"github.com/seemsod1/api-project/internal/logger"
	"github.com/seemsod1/api-project/internal/timezone"
)

const (
	TimeToSend   = 9 // 9 AM to send emails
	MinuteToSend = 1 // send at *:01 AM
)

type EmailNotifier struct {
	Subscriber    Subscriber
	Event         Event
	Scheduler     Scheduler
	RateService   RateService
	Logger        *logger.Logger
	KafkaProducer *kafka.Writer
}

type (
	EmailSender interface {
		Send(e *email.Email) error
	}

	RateService interface {
		GetRate(ctx context.Context, base, target string) (float64, error)
	}

	Scheduler interface {
		Start()
		AddEverydayJob(task func(), minute int) error
	}
	Subscriber interface {
		GetSubscribersWithTimezone(timezoneDiff int) ([]string, error)
	}
	Event interface {
		AddToEvents([]emailStreamer.Event) error
	}
)

func NewEmailNotifier(subs Subscriber, event Event, sch Scheduler, rateService RateService,
	logg *logger.Logger, kafkaProducer *kafka.Writer,
) *EmailNotifier {
	return &EmailNotifier{
		Subscriber:    subs,
		Event:         event,
		Scheduler:     sch,
		RateService:   rateService,
		Logger:        logg,
		KafkaProducer: kafkaProducer,
	}
}

func (et *EmailNotifier) Start() error {
	et.Logger.Info("Starting mail sender")

	sch := et.Scheduler
	sch.Start()
	err := sch.AddEverydayJob(func() {
		et.Logger.Info("Sending emails")

		localTime := time.Now().Hour()
		timezoneDiff := timezone.GetTimezoneDiff(localTime, TimeToSend)
		if err := timezone.ValidateTimezoneDiff(timezoneDiff); err != nil {
			et.Logger.Errorf("Error validating timezone diff: %v\n", err)
			return
		}

		subs, er := et.Subscriber.GetSubscribersWithTimezone(timezoneDiff)
		if er != nil {
			et.Logger.Errorf("Error getting subscribers: %v\n", er)
			return
		}

		et.SendRate(subs)
		et.Logger.Info("Emails sent")
	}, MinuteToSend)
	if err != nil {
		return fmt.Errorf("adding everyday job: %w", err)
	}

	return nil
}

func (et *EmailNotifier) SendRate(recipients []string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rate, err := et.RateService.GetRate(ctx, "USD", "UAH")
	if err != nil {
		et.Logger.Errorf("Error getting rate: %v\n", err)
		return
	}
	msgText := fmt.Sprintf("Current rate: %.2f", rate)

	messages := make([]emailStreamer.Event, 0, len(recipients))
	for _, recipient := range recipients {
		data := emailStreamer.Data{
			Recipient: recipient,
			Message:   msgText,
		}
		serializedData, er := emailStreamer.SerializeData(data)
		if er != nil {
			et.Logger.Errorf("Error serializing data: %v\n", err)
			return
		}

		msg := emailStreamer.Event{
			Data: serializedData,
		}
		messages = append(messages, msg)
	}

	if err = et.Event.AddToEvents(messages); err != nil {
		et.Logger.Errorf("Error adding to outbox: %v\n", err)
		return
	}

	et.Logger.Info("All messages saved to outbox")
}
