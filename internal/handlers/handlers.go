package handlers

import (
	"github.com/seemsod1/api-project/internal/storage"
)

// Repo is the repository used by the handlers
var Repo *Repository

// Repository is the repository struct
type Repository struct {
	DB storage.DatabaseRepo
}

// NewRepo creates a new repository with GORM
func NewRepo(db storage.DatabaseRepo) *Repository {
	return &Repository{
		DB: db,
	}
}

// NewHandlers creates a new handlers
func NewHandlers(r *Repository) {
	Repo = r
}
