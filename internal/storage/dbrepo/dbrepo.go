package dbrepo

import (
	"github.com/seemsod1/api-project/internal/config"
	"github.com/seemsod1/api-project/internal/storage"
	"gorm.io/gorm"
)

// gormDBRepo is a struct that embeds the AppConfig and the gorm.DB
type gormDBRepo struct {
	App *config.AppConfig
	DB  *gorm.DB
}

// NewGormRepo creates a new gormDBRepo
func NewGormRepo(conn *gorm.DB, a *config.AppConfig) storage.DatabaseRepo {
	return &gormDBRepo{
		App: a,
		DB:  conn,
	}
}
