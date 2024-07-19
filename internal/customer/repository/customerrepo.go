package customerrepo

import (
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/seemsod1/api-project/internal/customer/models"
	"gorm.io/gorm"
)

var ErrorDuplicateCustomer = errors.New("customer already exists")

type CustomerRepo struct {
	DB *gorm.DB
}

func NewCustomerRepo(db *gorm.DB) (*CustomerRepo, error) {
	if err := db.AutoMigrate(&customersmodels.Customer{}); err != nil {
		return nil, err
	}
	return &CustomerRepo{
		DB: db,
	}, nil
}

func (c *CustomerRepo) CreateCustomer(cust customersmodels.Customer) (int, error) {
	err := c.DB.Create(&cust).Error
	duplicateEntryError := &pgconn.PgError{Code: "23505"}

	if err != nil {
		if errors.As(err, &duplicateEntryError) {
			return 0, ErrorDuplicateCustomer
		}
		return 0, fmt.Errorf("gorm adding subscriber: %w", err)
	}
	return int(cust.ID), nil
}

func (c *CustomerRepo) UpdateCustomerStatus(id int) error {
	return c.DB.Model(&customersmodels.Customer{}).Where("id = ?", id).Update("status", "subscribed").Error
}

func (c *CustomerRepo) DeleteCustomerByID(id int) error {
	return c.DB.Where("id = ?", id).Delete(&customersmodels.Customer{}).Error
}
