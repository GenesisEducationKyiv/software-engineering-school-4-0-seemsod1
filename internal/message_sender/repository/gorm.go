package messagesenderrepo

import (
	"fmt"
	"time"

	messagesender "github.com/seemsod1/api-project/internal/message_sender"
	"gorm.io/gorm"
)

type EventDBRepo struct {
	DB *gorm.DB
}

func NewEventDBRepo(db *gorm.DB) (*EventDBRepo, error) {
	if err := db.AutoMigrate(&messagesender.EventProcessed{}); err != nil {
		return nil, fmt.Errorf("gorm migrating event processed: %w", err)
	}

	return &EventDBRepo{
		DB: db,
	}, nil
}

func (m *EventDBRepo) ConsumeEvent(event messagesender.EventProcessed) error {
	event.ProcessedAt = time.Now()
	if err := m.DB.Create(&event).Error; err != nil {
		return fmt.Errorf("gorm consuming event: %w", err)
	}
	return nil
}

func (m *EventDBRepo) CheckEventProcessed(id int) (bool, error) {
	var count int64
	err := m.DB.Model(&messagesender.EventProcessed{}).Where("id = ?", id).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("gorm checking event existence: %w", err)
	}
	return count > 0, nil
}
