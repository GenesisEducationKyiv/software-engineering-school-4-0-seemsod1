package notifierrepo

import (
	"context"
	"fmt"
	"time"

	"github.com/seemsod1/api-project/internal/notifier"
	"gorm.io/gorm"
)

type EventDBRepo struct {
	DB *gorm.DB
}

func NewEventDBRepo(db *gorm.DB) (*EventDBRepo, error) {
	if err := db.AutoMigrate(&notifier.Event{}); err != nil {
		return nil, fmt.Errorf("gorm migrating event: %w", err)
	}
	return &EventDBRepo{
		DB: db,
	}, nil
}

// AddToEvents adds a new message to the outbox
func (m *EventDBRepo) AddToEvents(messages []notifier.Event) error {
	tx := m.DB.Begin()
	if tx.Error != nil {
		return fmt.Errorf("gorm starting transaction: %w", tx.Error)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, msg := range messages {
		msgCopy := msg
		msgCopy.CreatedAt = time.Now()
		if err := tx.Create(&msgCopy).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("gorm adding message to outbox: %w", err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("gorm committing transaction: %w", err)
	}
	return nil
}

func (m *EventDBRepo) GetEvents(offset uint, limit int) ([]notifier.Event, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var messages []notifier.Event
	err := m.DB.WithContext(ctx).Where("id > ?", offset).Limit(limit).Find(&messages).Error
	if err != nil {
		return nil, fmt.Errorf("gorm getting outbox messages: %w", err)
	}
	return messages, nil
}
