package dbrepo

import (
	"github.com/seemsod1/api-project/internal/config"
	"github.com/seemsod1/api-project/internal/storage"
	"gorm.io/gorm"
)

type gormDBRepo struct {
	App *config.AppConfig
	DB  *gorm.DB
}

func NewGormRepo(conn *gorm.DB, a *config.AppConfig) storage.DatabaseRepo {
	return &gormDBRepo{
		App: a,
		DB:  conn,
	}
}
