package dbrepo

import (
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/seemsod1/api-project/internal/models"
)

var (
	ErrorDuplicateSubscription   = errors.New("subscription already exists")
	ErrorNonExistentSubscription = errors.New("subscription does not exist")
)

// AddSubscriber adds a new subscriber to the database
func (m *gormDBRepo) AddSubscriber(subscriber models.Subscriber) error {
	err := m.DB.Create(&subscriber).Error

	duplicateEntryError := &pgconn.PgError{Code: "23505"}

	if err != nil {
		if errors.As(err, &duplicateEntryError) {
			return ErrorDuplicateSubscription
		}
		return fmt.Errorf("gorm adding subscriber: %w", err)
	}
	return nil
}

// RemoveSubscriber removes a subscriber from the database
func (m *gormDBRepo) RemoveSubscriber(email string) error {
	res := m.DB.Where("email = ?", email).Delete(&models.Subscriber{})
	if res.Error != nil {
		return fmt.Errorf("gorm deleting subscriber: %w", res.Error)
	}

	if res.RowsAffected == 0 {
		return ErrorNonExistentSubscription
	}

	return nil
}

// GetSubscribersWithTimezone returns all subscribers from the database
func (m *gormDBRepo) GetSubscribersWithTimezone(timezone int) ([]string, error) {
	var emails []string
	err := m.DB.Model(&models.Subscriber{}).Where("timezone = ?", timezone).Pluck("email", &emails).Error
	if err != nil {
		return nil, fmt.Errorf("gorm getting subscribers: %w", err)
	}

	return emails, nil
}

// GetSubscribers returns all subscribers from the database
func (m *gormDBRepo) GetSubscribers() ([]string, error) {
	var emails []string
	err := m.DB.Model(&models.Subscriber{}).Pluck("email", &emails).Error
	if err != nil {
		return nil, fmt.Errorf("gorm getting subscribers: %w", err)
	}

	return emails, nil
}
