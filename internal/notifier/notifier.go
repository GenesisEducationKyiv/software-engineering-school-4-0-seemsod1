package notifier

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/seemsod1/api-project/pkg/notifier"

	"github.com/seemsod1/api-project/pkg/timezone"

	"github.com/seemsod1/api-project/pkg/logger"
)

const (
	timeToSend   = 9 // 9 AM to send emails
	minuteToSend = 1 // send at *:01 AM
	serviceName  = "notifier"
)

type EmailNotifier struct {
	Subscriber  subscriberRepo
	Event       eventRepo
	Scheduler   scheduler
	RateService rateService
	Logger      *logger.Logger
}

type (
	rateService interface {
		GetRate(ctx context.Context, base, target string) (float64, error)
	}

	scheduler interface {
		Start()
		AddEverydayJob(task func(), minute int) error
	}
	subscriberRepo interface {
		GetSubscribersWithTimezone(timezoneDiff int) ([]string, error)
	}
	eventRepo interface {
		AddToEvents([]notifier.Event) error
	}
)

func NewEmailNotifier(subs subscriberRepo, eventRepo eventRepo, sch scheduler, rateService rateService,
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
	ctx := context.Background()
	ctx = context.WithValue(ctx, logger.ServiceNameKey, serviceName)

	err := sch.AddEverydayJob(func() {
		et.Logger.WithContext(ctx).Info("Sending emails")

		localTime := time.Now().Hour()
		timezoneDiff := timezone.GetTimezoneDiff(localTime, timeToSend)
		if err := timezone.ValidateTimezoneDiff(timezoneDiff); err != nil {
			et.Logger.WithContext(ctx).Error("Error validating timezone diff", zap.Error(err))
			return
		}

		et.Logger.Debug("getting subscribers")
		subs, er := et.Subscriber.GetSubscribersWithTimezone(timezoneDiff)
		if er != nil {
			et.Logger.WithContext(ctx).Error("Error getting subscribers", zap.Error(er))
			return
		}

		et.SendRate(ctx, subs)
		et.Logger.WithContext(ctx).Info("Emails sent")
	}, minuteToSend)
	if err != nil {
		return fmt.Errorf("adding everyday job: %w", err)
	}

	return nil
}

func (et *EmailNotifier) SendRate(cntx context.Context, recipients []string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	et.Logger.WithContext(cntx).Debug("getting rate from providers")
	rate, err := et.RateService.GetRate(ctx, "USD", "UAH")
	if err != nil {
		et.Logger.WithContext(cntx).Error("Error getting rate", zap.Error(err))
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
			et.Logger.WithContext(cntx).Error("Error serializing data", zap.Error(er))
			return
		}

		msg := notifier.Event{
			Data: serializedData,
		}
		messages = append(messages, msg)
	}

	if err = et.Event.AddToEvents(messages); err != nil {
		et.Logger.WithContext(cntx).Error("Error adding to events list", zap.Error(err))
		return
	}

	et.Logger.WithContext(cntx).Info("All messages saved to outbox")
}

// serializeData converts a Message struct to a JSON string, excluding ID, CreatedAt and SentAt
func serializeData(data Data) (string, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to serialize message: %w", err)
	}
	return string(jsonData), nil
}
