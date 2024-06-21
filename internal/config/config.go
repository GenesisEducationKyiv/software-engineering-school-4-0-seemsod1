package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

// AppConfig is a struct that holds the configuration of the app
type AppConfig struct {
	Prod bool `required:"true"`
}

// NewAppConfig creates a new AppConfig
func NewAppConfig() (*AppConfig, error) {
	var appConfig AppConfig
	err := envconfig.Process("app", &appConfig)
	if err != nil {
		return nil, fmt.Errorf("processing app config: %w", err)
	}
	return &appConfig, nil
}
