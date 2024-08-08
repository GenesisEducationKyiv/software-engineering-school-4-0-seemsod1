package handlers

import (
	"github.com/seemsod1/api-project/pkg/logger"
)

// Handlers is a struct that contains all the services
type Handlers struct {
	Customer    customer
	Subscriber  subscriber
	RateService rateService
	Logger      *logger.Logger
}

// NewHandlers creates handlers
func NewHandlers(c customer, subs subscriber, rateService rateService, log *logger.Logger) *Handlers {
	return &Handlers{
		Customer:    c,
		Subscriber:  subs,
		RateService: rateService,
		Logger:      log,
	}
}
