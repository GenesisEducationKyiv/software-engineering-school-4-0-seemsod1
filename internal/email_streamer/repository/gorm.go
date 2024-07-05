package streamerrepo

import (
	"errors"
	"fmt"
	"time"

	emailStreamer "github.com/seemsod1/api-project/internal/email_streamer"
	"gorm.io/gorm"
)

type Streamer struct {
	DB *gorm.DB
}

// NewStreamerRepo creates a new StreamerRepo
func NewStreamerRepo(db *gorm.DB) (*Streamer, error) {
	if err := db.AutoMigrate(&emailStreamer.Event{}); err != nil {
		return nil, fmt.Errorf("gorm migrating emailStreamer.Event: %w", err)
	}
	if err := db.AutoMigrate(&emailStreamer.Streamer{}); err != nil {
		return nil, fmt.Errorf("gorm migrating emailStreamer.Event: %w", err)
	}
	if err := db.AutoMigrate(&emailStreamer.EventProcessed{}); err != nil {
		return nil, fmt.Errorf("gorm migrating emailStreamer.Event: %w", err)
	}

	return &Streamer{DB: db}, nil
}

// AddToEvents adds a new message to the outbox
func (m *Streamer) AddToEvents(messages []emailStreamer.Event) error {
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

func (m *Streamer) GetOutboxEvents(offset uint, limit int) ([]emailStreamer.Event, error) {
	var messages []emailStreamer.Event
	err := m.DB.Where("id > ?", offset).Limit(limit).Find(&messages).Error
	if err != nil {
		return nil, fmt.Errorf("gorm getting outbox messages: %w", err)
	}
	return messages, nil
}

func (m *Streamer) ChangeOffset(msg emailStreamer.Streamer) error {
	var existing emailStreamer.Streamer
	err := m.DB.Where("topic = ? AND partition = ?", msg.Topic, msg.Partition).First(&existing).Error

	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		if err = m.DB.Create(&msg).Error; err != nil {
			return fmt.Errorf("gorm adding message to stream: %w", err)
		}
	case err == nil:
		existing.LastOffset = msg.LastOffset
		if err = m.DB.Save(&existing).Error; err != nil {
			return fmt.Errorf("gorm updating message to stream: %w", err)
		}
	default:
		return fmt.Errorf("gorm retrieving existing record: %w", err)
	}
	return nil
}

func (m *Streamer) ConsumeEvent(event emailStreamer.EventProcessed) error {
	event.ProcessedAt = time.Now()
	if err := m.DB.Create(&event).Error; err != nil {
		return fmt.Errorf("gorm consuming event: %w", err)
	}
	return nil
}

func (m *Streamer) CheckEventProcessed(id int) (bool, error) {
	var count int64
	err := m.DB.Model(&emailStreamer.EventProcessed{}).Where("id = ?", id).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("gorm checking event existence: %w", err)
	}
	return count > 0, nil
}

func (m *Streamer) GetOffset(topic string, partition int) (uint, error) {
	var offset uint
	err := m.DB.Model(&emailStreamer.Streamer{}).
		Where("topic = ? AND partition = ?", topic, partition).
		Pluck("last_offset", &offset).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, nil
		}
		return 0, fmt.Errorf("error retrieving the latest offset: %w", err)
	}
	return offset, nil
}
