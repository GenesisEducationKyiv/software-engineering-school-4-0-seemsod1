package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

// AppConfig is a struct that holds the configuration of the app
type AppConfig struct{}

// NewAppConfig creates a new AppConfig
func NewAppConfig() *AppConfig {
	var appConfig AppConfig
	if err := envconfig.Process("app", &appConfig); err != nil {
		log.Fatal(err)
	}
	return &appConfig
}
