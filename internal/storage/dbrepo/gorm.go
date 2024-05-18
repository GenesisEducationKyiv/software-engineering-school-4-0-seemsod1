package dbrepo

import (
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
	customerrors "github.com/seemsod1/api-project/internal/errors"
	"github.com/seemsod1/api-project/internal/models"
)

// AddSubscriber adds a new subscriber to the database
func (m *gormDBRepo) AddSubscriber(subscriber models.Subscriber) error {
	err := m.DB.Create(&subscriber).Error

	var duplicateEntryError = &pgconn.PgError{Code: "23505"}

	if err != nil {
		if errors.As(err, &duplicateEntryError) {
			return customerrors.DuplicatedKey
		}
		return err
	}
	return nil
}

// GetSubscribers returns all subscribers from the database
func (m *gormDBRepo) GetSubscribers() ([]string, error) {
	var subscribers []models.Subscriber
	err := m.DB.Find(&subscribers).Error
	if err != nil {
		return nil, err
	}

	var emails []string
	for _, s := range subscribers {
		emails = append(emails, s.Email)
	}
	return emails, nil
}
