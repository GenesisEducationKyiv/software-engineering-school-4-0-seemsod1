package notifier

import (
	"fmt"
	"log"
	"net/smtp"
	"sync"

	"github.com/jordan-wright/email"
	"github.com/kelseyhightower/envconfig"
	"github.com/seemsod1/api-project/internal/driver"
	"github.com/seemsod1/api-project/internal/rateapi"
	"github.com/seemsod1/api-project/internal/scheduler"
	"github.com/seemsod1/api-project/internal/storage"
	"github.com/seemsod1/api-project/internal/storage/dbrepo"
	"github.com/seemsod1/api-project/internal/timezone"
)

// EmailNotifier is a struct that contains the app configuration and the database repository.
type EmailNotifier struct {
	DB storage.DatabaseRepo
}

type EmailNotifierConfig struct {
	Host     string
	Port     string
	From     string
	Password string
}

const (
	numWorkers   = 4
	bufferSize   = 100
	timeToSend   = 9 // 9 AM to send emails
	minuteToSend = 1 // send at *:01 AM
)

// NewEmailNotifierWithGORM creates a new email notifier with GORM.
func NewEmailNotifierWithGORM(db *driver.GORMDriver) *EmailNotifier {
	return &EmailNotifier{DB: dbrepo.NewGormRepo(db.SQL)}
}

func NewEmailNotifierConfig() EmailNotifierConfig {
	var cfg EmailNotifierConfig
	if err := envconfig.Process("mailer", &cfg); err != nil {
		log.Fatal(err)
	}
	return cfg
}

func (c *EmailNotifierConfig) Validate() bool {
	if c.Host == "" || c.Port == "" || c.From == "" || c.Password == "" {
		return false
	}
	return true
}

// Start starts the mail sender.
func (et *EmailNotifier) Start() error {
	log.Println("Starting mail sender")
	cfg := NewEmailNotifierConfig()

	if !cfg.Validate() {
		return fmt.Errorf("invalid mail config")
	}

	sch := scheduler.NewGoCronScheduler()
	sch.Start()
	err := sch.AddEverydayJob(func() {
		log.Println("Sending emails")

		timezoneDiff := timezone.GetTimezoneDiff(timeToSend)

		subs, err := et.DB.GetSubscribers(timezoneDiff)
		if err != nil {
			log.Printf("Error getting subscribers: %v\n", err)
			return
		}
		et.sendEmails(cfg, subs)
		log.Println("Emails sent")
	}, minuteToSend)

	return err
}

// initPool initializes the email pool.
func (et *EmailNotifier) initPool(cfg EmailNotifierConfig) (*email.Pool, chan<- *email.Email, *sync.WaitGroup, error) {
	ch := make(chan *email.Email, bufferSize)
	var wg sync.WaitGroup

	p, err := email.NewPool(
		fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		numWorkers,
		smtp.PlainAuth("", cfg.From, cfg.Password, cfg.Host),
	)
	if err != nil {
		return nil, nil, nil, err
	}

	for i := 0; i < numWorkers; i++ {
		go func() {
			for e := range ch {
				err := p.Send(e, -1)
				if err != nil {
					log.Printf("Error sending email: %v\n", err)
				}
				wg.Done()
			}
		}()
	}

	return p, ch, &wg, nil
}

// sendEmail sends an email.
func (et *EmailNotifier) sendEmail(ch chan<- *email.Email, wg *sync.WaitGroup, e *email.Email) {
	wg.Add(1)
	ch <- e
}

// sendEmails sends emails to the recipients.
func (et *EmailNotifier) sendEmails(cfg EmailNotifierConfig, recipients []string) {
	p, ch, wg, err := et.initPool(cfg)
	if err != nil {
		log.Printf("Error initializing email pool: %v\n", err)
		return
	}

	provider := rateapi.NewCoinbaseProvider()

	price, err := provider.GetUsdToUahRate()
	if err != nil {
		log.Printf("Error getting rate: %v\n", err)
		return
	}
	msgText := []byte("Current rate: " + fmt.Sprintf("%.2f", price))

	for _, recipient := range recipients {
		e := email.NewEmail()
		e.From = cfg.From
		e.To = []string{recipient}
		e.Subject = "Currency rate notification: USD to UAH"
		e.Text = msgText
		et.sendEmail(ch, wg, e)
	}

	wg.Wait()
	p.Close()
	close(ch)
}
