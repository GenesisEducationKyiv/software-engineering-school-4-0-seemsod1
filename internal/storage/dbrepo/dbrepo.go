package dbrepo

import (
	"github.com/seemsod1/api-project/internal/config"
	"github.com/seemsod1/api-project/internal/storage"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// gormDBRepo is a struct that embeds the AppConfig and the gorm.DB
type gormDBRepo struct {
	App *config.AppConfig
	DB  *gorm.DB
}

type mockDB struct {
	mock.Mock
}

// NewGormRepo creates a new gormDBRepo
func NewGormRepo(conn *gorm.DB, a *config.AppConfig) storage.DatabaseRepo {
	return &gormDBRepo{
		App: a,
		DB:  conn,
	}
}

func NewMockDB() *mockDB {
	return &mockDB{}
}
