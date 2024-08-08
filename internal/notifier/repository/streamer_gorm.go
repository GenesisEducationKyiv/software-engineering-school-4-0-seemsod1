package notifierrepo

import (
	"errors"
	"fmt"

	emailStreamer "github.com/seemsod1/api-project/pkg/email_streamer"
	"gorm.io/gorm"
)

type Streamer struct {
	DB *gorm.DB
}

// NewStreamerRepo creates a new StreamerRepo
func NewStreamerRepo(db *gorm.DB) (*Streamer, error) {
	if err := db.AutoMigrate(&emailStreamer.Streamer{}); err != nil {
		return nil, fmt.Errorf("gorm migrating streamer: %w", err)
	}
	return &Streamer{
		DB: db,
	}, nil
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
