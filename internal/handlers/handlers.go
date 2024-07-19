package handlers

import (
	"github.com/seemsod1/api-project/pkg/logger"
)

// Repo is the repository used by the handlers
var Repo *Repository

// Repository is the repository struct
type Repository struct {
	Customer    Customer
	Subscriber  Subscriber
	RateService RateService
	Logger      *logger.Logger
}

// NewRepo creates a new repository with GORM
func NewRepo(c Customer, subs Subscriber, rateService RateService, log *logger.Logger) *Repository {
	return &Repository{
		Customer:    c,
		Subscriber:  subs,
		RateService: rateService,
		Logger:      log,
	}
}

// NewHandlers creates a new handlers
func NewHandlers(r *Repository) {
	Repo = r
}
