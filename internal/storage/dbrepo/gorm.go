package dbrepo

import (
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	customerrors "github.com/seemsod1/api-project/internal/errors"
	"github.com/seemsod1/api-project/internal/models"
)

// AddSubscriber adds a new subscriber to the database
func (m *gormDBRepo) AddSubscriber(subscriber models.Subscriber) error {
	err := m.DB.Create(&subscriber).Error

	duplicateEntryError := &pgconn.PgError{Code: "23505"}

	if err != nil {
		if errors.As(err, &duplicateEntryError) {
			return customerrors.ErrDuplicatedKey
		}
		return fmt.Errorf("gorm adding subscriber: %w", err)
	}
	return nil
}

// GetSubscribers returns all subscribers from the database
func (m *gormDBRepo) GetSubscribers(timezone int) ([]string, error) {
	var emails []string
	err := m.DB.Model(&models.Subscriber{}).Where("timezone = ?", timezone).Pluck("email", &emails).Error
	if err != nil {
		return nil, fmt.Errorf("gorm getting subscribers: %w", err)
	}

	return emails, nil
}
