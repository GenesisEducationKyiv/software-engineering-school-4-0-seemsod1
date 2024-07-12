package subscriberrepo

import (
	"context"
	"errors"
	"fmt"
	"time"

	subscribermodels "github.com/seemsod1/api-project/internal/subscriber/models"

	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrorDuplicateSubscription   = errors.New("subscription already exists")
	ErrorNonExistentSubscription = errors.New("subscription does not exist")
)

// AddSubscriber adds a new subscriber to the database
func (m *SubscriberDBRepo) AddSubscriber(subs subscribermodels.Subscriber) error {
	err := m.DB.Create(&subs).Error

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
func (m *SubscriberDBRepo) RemoveSubscriber(email string) error {
	res := m.DB.Where("email = ?", email).Delete(&subscribermodels.Subscriber{})
	if res.Error != nil {
		return fmt.Errorf("gorm deleting subscriber: %w", res.Error)
	}

	if res.RowsAffected == 0 {
		return ErrorNonExistentSubscription
	}

	return nil
}

// GetSubscribersWithTimezone returns all subscribers from the database
func (m *SubscriberDBRepo) GetSubscribersWithTimezone(timezone int) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var emails []string
	err := m.DB.WithContext(ctx).Model(&subscribermodels.Subscriber{}).Where("timezone = ?", timezone).Pluck("email", &emails).Error
	if err != nil {
		return nil, fmt.Errorf("gorm getting subscribers: %w", err)
	}

	return emails, nil
}

// GetSubscribers returns all subscribers from the database
func (m *SubscriberDBRepo) GetSubscribers() ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var emails []string
	err := m.DB.WithContext(ctx).Model(&subscribermodels.Subscriber{}).Pluck("email", &emails).Error
	if err != nil {
		return nil, fmt.Errorf("gorm getting subscribers: %w", err)
	}

	return emails, nil
}
