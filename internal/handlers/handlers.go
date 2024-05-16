package handlers

import (
	"github.com/seemsod1/api-project/internal/config"
	"github.com/seemsod1/api-project/internal/driver"
	"github.com/seemsod1/api-project/internal/storage"
	"github.com/seemsod1/api-project/internal/storage/dbrepo"
)

var Repo *Repository

type Repository struct {
	DB  storage.DatabaseRepo
	App *config.AppConfig
}

func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewGormRepo(db.SQL, a),
	}
}

func NewHandlers(r *Repository) {
	Repo = r
}
