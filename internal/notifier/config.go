package notifier

import "github.com/kelseyhightower/envconfig"

type EmailNotifierConfig struct {
	Host     string `required:"true"`
	Port     string `required:"true"`
	From     string `required:"true"`
	Password string `required:"true"`
}

func NewEmailNotifierConfig() (EmailNotifierConfig, error) {
	var cfg EmailNotifierConfig
	if err := envconfig.Process("mailer", &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}

func (c *EmailNotifierConfig) Validate() bool {
	return !(c.Host == "" || c.Port == "" || c.From == "" || c.Password == "")
}
