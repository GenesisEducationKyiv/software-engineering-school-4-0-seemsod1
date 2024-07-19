package subscriberrepo

import (
	"fmt"

	subscribermodels "github.com/seemsod1/api-project/internal/subscriber/models"
	"gorm.io/gorm"
)

type SubscriberDBRepo struct {
	DB *gorm.DB
}

func NewSubscriberDBRepo(db *gorm.DB) (*SubscriberDBRepo, error) {
	if err := db.AutoMigrate(&subscribermodels.Subscriber{}); err != nil {
		return nil, fmt.Errorf("gorm migrating event: %w", err)
	}
	return &SubscriberDBRepo{
		DB: db,
	}, nil
}
