package customer

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	customersmodels "github.com/seemsod1/api-project/internal/customer/models"
	customerrepo "github.com/seemsod1/api-project/internal/customer/repository"
	"github.com/seemsod1/api-project/pkg/logger"
	"github.com/segmentio/kafka-go"
)

type SagaCoordinator struct {
	CustomerRepo Database
	Producer     *kafka.Writer
	Consumer     *kafka.Reader
	Logger       *logger.Logger
}

type Database interface {
	CreateCustomer(customer customersmodels.Customer) (int, error)
	UpdateCustomerStatus(id int) error
	DeleteCustomerByID(id int) error
}

func NewSagaCoordinator(customer Database, producer *kafka.Writer, consumer *kafka.Reader, logg *logger.Logger) *SagaCoordinator {
	return &SagaCoordinator{
		CustomerRepo: customer,
		Producer:     producer,
		Consumer:     consumer,
		Logger:       logg,
	}
}

func (c *SagaCoordinator) StartTransaction(email string, timezone int) error {
	customer := customersmodels.Customer{
		Email:    email,
		Timezone: timezone,
	}

	id, err := c.CustomerRepo.CreateCustomer(customer)
	if err != nil {
		if errors.Is(err, customerrepo.ErrorDuplicateCustomer) {
			return customerrepo.ErrorDuplicateCustomer
		}
		c.Logger.Error("failed to create customer")
		return fmt.Errorf("creating customer: %w", err)
	}
	data := customersmodels.CommandData{
		Command: "subscribe_by_email",
		Payload: customersmodels.Payload{
			Email:    customer.Email,
			Timezone: customer.Timezone,
		},
	}

	serializedData, er := customersmodels.SerializeCommandData(data)
	if er != nil {
		_ = c.CustomerRepo.DeleteCustomerByID(id)
		c.Logger.Error("failed to serialize data")
		return fmt.Errorf("serializing data: %w", er)
	}

	serializedID := fmt.Sprintf("%d", id)

	msg := kafka.Message{
		Key:   []byte(serializedID),
		Value: []byte(serializedData),
	}
	if err = c.Producer.WriteMessages(context.Background(), msg); err != nil {
		_ = c.CustomerRepo.DeleteCustomerByID(id)
		c.Logger.Error("failed to send message")
		return fmt.Errorf("sending message: %w", err)
	}
	return nil
}

func (c *SagaCoordinator) StartReceivingMessages(ctx context.Context) {
	for {
		m, err := c.Consumer.FetchMessage(ctx)
		if err != nil {
			c.Logger.Error("failed to read message")
			continue
		}

		msg, err := customersmodels.DeserializeReplyData(string(m.Value))
		if err != nil {
			c.Logger.Error("failed to deserialize data")
			continue
		}

		switch msg.Command {
		case "subscribe_by_email":
			switch msg.Reply {
			case "success":
				id, er := strconv.Atoi(string(m.Key))
				if er != nil {
					c.Logger.Error("failed to convert id")
					continue
				}
				if err = c.CustomerRepo.UpdateCustomerStatus(id); err != nil {
					c.Logger.Error("failed to update customer status")
					continue
				}
			default:
				c.Logger.Error("unknown reply")
			}

		default:
			c.Logger.Error("unknown command")
		}

		if err = c.Consumer.CommitMessages(ctx, m); err != nil {
			c.Logger.Error("failed to commit message")
		}
	}
}
