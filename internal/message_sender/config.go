package messagesender

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type EmailSenderConfig struct {
	Host     string `required:"true"`
	Port     string `required:"true"`
	From     string `required:"true"`
	Password string `required:"true"`
}

func NewEmailSenderConfig() (EmailSenderConfig, error) {
	var cfg EmailSenderConfig
	if err := envconfig.Process("mailer", &cfg); err != nil {
		return cfg, fmt.Errorf("processing mailer config: %w", err)
	}
	return cfg, nil
}

func (c *EmailSenderConfig) Validate() bool {
	return !(c.Host == "" || c.Port == "" || c.From == "" || c.Password == "")
}
