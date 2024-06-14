package handlers

import (
	"github.com/seemsod1/api-project/internal/driver"
	"github.com/seemsod1/api-project/internal/storage"
	"github.com/seemsod1/api-project/internal/storage/dbrepo"
)

// Repo is the repository used by the handlers
var Repo *Repository

// Repository is the repository struct
type Repository struct {
	DB storage.DatabaseRepo
}

// NewRepoWithGORM creates a new repository with GORM
func NewRepoWithGORM(db *driver.GORMDriver) *Repository {
	return &Repository{
		DB: dbrepo.NewGormRepo(db.SQL),
	}
}

// NewHandlers creates a new handlers
func NewHandlers(r *Repository) {
	Repo = r
}
