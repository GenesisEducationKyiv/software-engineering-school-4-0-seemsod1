package notifier

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/seemsod1/api-project/internal/logger"
	"go.uber.org/zap"

	"github.com/jordan-wright/email"
	"github.com/seemsod1/api-project/internal/timezone"
)

const (
	numWorkers   = 4
	bufferSize   = 100
	TimeToSend   = 9 // 9 AM to send emails
	MinuteToSend = 1 // send at *:01 AM
)

type EmailNotifier struct {
	Subscriber  Subscriber
	Scheduler   Scheduler
	RateService RateService
	EmailSender EmailSender
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
	Subscriber interface {
		GetSubscribers(timezoneDiff int) ([]string, error)
	}
)

func NewEmailNotifier(subs Subscriber, sch Scheduler, rateService RateService,
	emailSender EmailSender, logg *logger.Logger,
) *EmailNotifier {
	return &EmailNotifier{
		Subscriber:  subs,
		Scheduler:   sch,
		RateService: rateService,
		EmailSender: emailSender,
		Logger:      logg,
	}
}

func (et *EmailNotifier) Start() error {
	et.Logger.Info("Starting mail sender")
	cfg, err := NewEmailNotifierConfig()
	if err != nil {
		return fmt.Errorf("creating mail sender config: %w", err)
	}

	if !cfg.Validate() {
		return fmt.Errorf("invalid mail config")
	}

	sch := et.Scheduler
	sch.Start()
	err = sch.AddEverydayJob(func() {
		et.Logger.Info("Sending emails")

		localTime := time.Now().Hour()
		timezoneDiff := timezone.GetTimezoneDiff(localTime, TimeToSend)
		if err = timezone.ValidateTimezoneDiff(timezoneDiff); err != nil {
			et.Logger.Errorf("Error validating timezone diff: %v\n", err)
			return
		}

		subs, er := et.Subscriber.GetSubscribers(timezoneDiff)
		if er != nil {
			et.Logger.Errorf("Error getting subscribers: %v\n", er)
			return
		}

		et.sendEmails(cfg, subs)
		et.Logger.Info("Emails sent")
	}, MinuteToSend)
	if err != nil {
		return fmt.Errorf("adding everyday job: %w", err)
	}

	return nil
}

func (et *EmailNotifier) sendEmails(cfg EmailNotifierConfig, recipients []string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rate, err := et.RateService.GetRate(ctx, "USD", "UAH")
	if err != nil {
		et.Logger.Errorf("Error getting rate: %v\n", err)
		return
	}
	msgText := []byte("Current rate: " + fmt.Sprintf("%.2f", rate))

	ch := make(chan *email.Email, bufferSize)
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		go func() {
			for e := range ch {
				if err := et.EmailSender.Send(e); err != nil {
					et.Logger.Error("Error sending email", zap.String("recipient", e.To[0]), zap.Error(err))
				}
				wg.Done()
			}
		}()
	}

	for _, recipient := range recipients {
		e := email.NewEmail()
		e.From = cfg.From
		e.To = []string{recipient}
		e.Subject = "Currency rate notification: USD to UAH"
		e.Text = msgText
		wg.Add(1)
		ch <- e
	}

	wg.Wait()
	close(ch)
}
