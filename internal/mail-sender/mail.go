package mailsender

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
	"strings"
	"sync"

	"github.com/go-co-op/gocron/v2"
	"github.com/jordan-wright/email"
	"github.com/seemsod1/api-project/internal/api"
	"github.com/seemsod1/api-project/internal/config"
	"github.com/seemsod1/api-project/internal/driver"
	"github.com/seemsod1/api-project/internal/helpers"
	"github.com/seemsod1/api-project/internal/storage"
	"github.com/seemsod1/api-project/internal/storage/dbrepo"
)

// MailSender is a struct that contains the app configuration and the database repository.
type MailSender struct {
	App    *config.AppConfig
	DB     storage.DatabaseRepo
	Config *Config
}

type Config struct {
	Host     string
	Port     string
	From     string
	Password string
}

const (
	numWorkers = 4
	bufferSize = 100
	timeToSend = 22 // 9 AM to send emails
)

// NewMailSender creates a new MailSender struct.
func NewMailSender(a *config.AppConfig, db *driver.DB) *MailSender {
	cfg := NewConfig()
	return &MailSender{DB: dbrepo.NewGormRepo(db.SQL, a), Config: cfg}
}

func NewConfig() *Config {
	host, port, from, pass := initConfigFromEnv()
	return &Config{Host: host, Port: port, From: from, Password: pass}
}

func (c *Config) Validate() bool {
	if c.Host == "" || c.Port == "" || c.From == "" || c.Password == "" {
		return false
	}
	return true
}

// Start starts the mail sender.
func (ms *MailSender) Start() error {
	log.Println("Starting mail sender")
	if !ms.Config.Validate() {
		return fmt.Errorf("invalid mail config")
	}

	s, err := gocron.NewScheduler()
	if err != nil {
		return err
	}
	s.Start()
	_, _ = s.NewJob(
		gocron.CronJob(
			"17 * * * *", // every hour at minute 1
			false,
		),
		gocron.NewTask(func() {
			log.Println("Sending emails")

			timezoneDiff := helpers.GetTimezoneDiff(timeToSend)

			subs, err := ms.DB.GetSubscribers(timezoneDiff)
			if err != nil {
				log.Printf("Error getting subscribers: %v\n", err)
				return
			}
			sendEmails(ms.Config, subs)
			log.Println("Emails sent")
			log.Println(subs)
		}),
	)

	return nil
}

// initPool initializes the email pool.
func initPool(cfg *Config) (*email.Pool, chan<- *email.Email, *sync.WaitGroup, error) {
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
func sendEmail(ch chan<- *email.Email, wg *sync.WaitGroup, e *email.Email) {
	wg.Add(1)
	ch <- e
}

// sendEmails sends emails to the recipients.
func sendEmails(cfg *Config, recipients []string) {
	p, ch, wg, err := initPool(cfg)
	if err != nil {
		log.Printf("Error initializing email pool: %v\n", err)
		return
	}

	price, err := api.GetUsdToUahRate()
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
		sendEmail(ch, wg, e)
	}

	wg.Wait()
	p.Close()
	close(ch)
}

func initConfigFromEnv() (hostName, hostPort, mailBox, appPassword string) {
	mailConfig := os.Getenv("MAILER_URL")
	if mailConfig == "" {
		return "", "", "", ""
	}

	parts := strings.Split(mailConfig, " ")
	if len(parts) != 4 {
		return "", "", "", ""
	}
	return parts[0], parts[1], parts[2], parts[3]
}
