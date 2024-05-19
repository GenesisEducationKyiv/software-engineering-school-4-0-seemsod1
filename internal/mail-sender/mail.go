package mailSender

import (
	"fmt"
	"github.com/go-co-op/gocron/v2"
	"github.com/jordan-wright/email"
	"github.com/seemsod1/api-project/internal/api"
	"github.com/seemsod1/api-project/internal/config"
	"github.com/seemsod1/api-project/internal/driver"
	"github.com/seemsod1/api-project/internal/helpers"
	"github.com/seemsod1/api-project/internal/storage"
	"github.com/seemsod1/api-project/internal/storage/dbrepo"
	"log"
	"net/smtp"
	"os"
	"sync"
)

// MailSender is a struct that contains the app configuration and the database repository.
type MailSender struct {
	App *config.AppConfig
	DB  storage.DatabaseRepo
}

const (
	numWorkers = 4
	bufferSize = 100
	timeToSend = 9 // 9 AM to send emails
)

// NewMailSender creates a new MailSender struct.
func NewMailSender(a *config.AppConfig, db *driver.DB) *MailSender {
	return &MailSender{DB: dbrepo.NewGormRepo(db.SQL, a)}
}

// Start starts the mail sender.
func (ms *MailSender) Start() error {
	log.Println("Starting mail sender")
	s, err := gocron.NewScheduler()
	if err != nil {
		return err
	}
	s.Start()
	_, _ = s.NewJob(
		gocron.CronJob(
			"1 * * * *", // every hour at minute 1
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
			sendEmails(subs)
			log.Println("Emails sent")
		}),
	)

	return nil
}

// initPool initializes the email pool.
func initPool() (*email.Pool, chan<- *email.Email, *sync.WaitGroup, error) {
	ch := make(chan *email.Email, bufferSize)
	var wg sync.WaitGroup

	p, err := email.NewPool(
		fmt.Sprintf("%s:%s", os.Getenv("MAIL_HOST"), os.Getenv("MAIL_PORT")),
		numWorkers,
		smtp.PlainAuth("", os.Getenv("MAIL_FROM"), os.Getenv("MAIL_PASS"), os.Getenv("MAIL_HOST")),
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
func sendEmails(recipients []string) {
	p, ch, wg, err := initPool()
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

	from := os.Getenv("MAIL_FROM")

	for _, recipient := range recipients {
		e := email.NewEmail()
		e.From = from
		e.To = []string{recipient}
		e.Subject = "Currency rate notification: USD to UAH"
		e.Text = msgText
		sendEmail(ch, wg, e)
	}

	wg.Wait()
	p.Close()
	close(ch)
}
