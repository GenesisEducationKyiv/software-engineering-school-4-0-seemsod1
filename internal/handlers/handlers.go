package handlers

import (
	"github.com/seemsod1/api-project/internal/storage"
)

// Repo is the repository used by the handlers
var Repo *Repository

// Repository is the repository struct
type Repository struct {
	DB          storage.DatabaseRepo
	RateService RateService
}

// NewRepo creates a new repository with GORM
func NewRepo(db storage.DatabaseRepo, rateService RateService) *Repository {
	return &Repository{
		DB:          db,
		RateService: rateService,
	}
}

// NewHandlers creates a new handlers
func NewHandlers(r *Repository) {
	Repo = r
}
