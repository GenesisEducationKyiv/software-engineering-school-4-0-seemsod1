package dbrepo

import (
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
	customerrors "github.com/seemsod1/api-project/internal/errors"
	"github.com/seemsod1/api-project/internal/models"
	"gorm.io/gorm"
)

func (m *gormDBRepo) Connection() *gorm.DB {
	return m.DB
}

func (m *gormDBRepo) AddSubscriber(email string) error {
	s := models.Subscriber{
		Email: email,
	}
	err := m.DB.Create(&s).Error

	var duplicateEntryError = &pgconn.PgError{Code: "23505"}

	if err != nil {
		if errors.As(err, &duplicateEntryError) {
			return customerrors.DuplicatedKey
		}
		return err
	}
	return nil
}
