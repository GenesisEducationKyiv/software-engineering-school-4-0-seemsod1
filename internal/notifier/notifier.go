package notifier

import (
	"context"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"time"

	"github.com/seemsod1/api-project/pkg/notifier"

	"github.com/seemsod1/api-project/pkg/timezone"

	"github.com/jordan-wright/email"
	"github.com/seemsod1/api-project/pkg/logger"
)

const (
	TimeToSend   = 9 // 9 AM to send emails
	MinuteToSend = 1 // send at *:01 AM
)

type EmailNotifier struct {
	Subscriber  SubscriberRepo
	Event       EventRepo
	Scheduler   Scheduler
	RateService RateService
	Logger      *logger.Logger
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
	SubscriberRepo interface {
		GetSubscribersWithTimezone(timezoneDiff int) ([]string, error)
	}
	EventRepo interface {
		AddToEvents([]notifier.Event) error
	}
)

func NewEmailNotifier(subs SubscriberRepo, eventRepo EventRepo, sch Scheduler, rateService RateService,
	logg *logger.Logger,
) *EmailNotifier {
	return &EmailNotifier{
		Subscriber:  subs,
		Event:       eventRepo,
		Scheduler:   sch,
		RateService: rateService,
		Logger:      logg,
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
			et.Logger.Error("Error validating timezone diff", zap.Error(err))
			return
		}

		subs, er := et.Subscriber.GetSubscribersWithTimezone(timezoneDiff)
		if er != nil {
			et.Logger.Error("Error getting subscribers", zap.Error(er))
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
		et.Logger.Error("Error getting rate", zap.Error(err))
		return
	}
	msgText := fmt.Sprintf("Current rate: %.2f", rate)

	messages := make([]notifier.Event, 0, len(recipients))
	for _, recipient := range recipients {
		data := Data{
			Recipient: recipient,
			Message:   msgText,
		}
		serializedData, er := serializeData(data)
		if er != nil {
			et.Logger.Error("Error serializing data", zap.Error(er))
			return
		}

		msg := notifier.Event{
			Data: serializedData,
		}
		messages = append(messages, msg)
	}

	if err = et.Event.AddToEvents(messages); err != nil {
		et.Logger.Error("Error adding to events list", zap.Error(err))
		return
	}

	et.Logger.Info("All messages saved to outbox")
}

// serializeData converts a Message struct to a JSON string, excluding ID, CreatedAt and SentAt
func serializeData(data Data) (string, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to serialize message: %w", err)
	}
	return string(jsonData), nil
}
