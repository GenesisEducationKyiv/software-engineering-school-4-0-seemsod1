package dbrepo

import (
	"github.com/seemsod1/api-project/internal/storage"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// gormDBRepo is a struct that embeds the AppConfig and the gorm.DB
type gormDBRepo struct {
	DB *gorm.DB
}

type MockDB struct {
	mock.Mock
}

// NewGormRepo creates a new gormDBRepo
func NewGormRepo(conn *gorm.DB) storage.DatabaseRepo {
	return &gormDBRepo{
		DB: conn,
	}
}

func NewMockDB() *MockDB {
	return &MockDB{}
}
