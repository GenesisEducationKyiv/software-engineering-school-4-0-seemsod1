package notifier

import (
	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"
)

type SMTPEmailSender struct {
	Pool *email.Pool
}

func NewSMTPEmailSender(cfg EmailNotifierConfig) (*SMTPEmailSender, error) {
	p, err := email.NewPool(
		fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		numWorkers,
		smtp.PlainAuth("", cfg.From, cfg.Password, cfg.Host),
	)
	if err != nil {
		return nil, fmt.Errorf("creating email pool: %w", err)
	}
	return &SMTPEmailSender{Pool: p}, nil
}

func (s *SMTPEmailSender) Send(e *email.Email) error {
	return s.Pool.Send(e, -1)
}
