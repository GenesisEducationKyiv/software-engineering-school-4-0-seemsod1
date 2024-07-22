package handlers

import (
	"github.com/seemsod1/api-project/pkg/logger"
)

// Handlers is a struct that contains all the services
type Handlers struct {
	Customer    Customer
	Subscriber  Subscriber
	RateService RateService
	Logger      *logger.Logger
}

// NewHandlers creates handlers
func NewHandlers(c Customer, subs Subscriber, rateService RateService, log *logger.Logger) *Handlers {
	return &Handlers{
		Customer:    c,
		Subscriber:  subs,
		RateService: rateService,
		Logger:      log,
	}
}
