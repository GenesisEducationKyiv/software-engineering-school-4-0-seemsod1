package customer

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/google/uuid"

	customersmodels "github.com/seemsod1/api-project/internal/customer/models"
	customerrepo "github.com/seemsod1/api-project/internal/customer/repository"
	"github.com/seemsod1/api-project/pkg/logger"
	"github.com/segmentio/kafka-go"
)

const serviceName = "customer"

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
	ctx := context.Background()
	traceID := uuid.New()
	ctx = context.WithValue(ctx, logger.TraceIDKey, traceID.String())
	ctx = context.WithValue(ctx, logger.ServiceNameKey, serviceName)

	customer := customersmodels.Customer{
		Email:    email,
		Timezone: timezone,
	}
	c.Logger.WithContext(ctx).Info("starting transaction")

	c.Logger.WithContext(ctx).Debug("creating customer")
	id, err := c.CustomerRepo.CreateCustomer(customer)
	if err != nil {
		if errors.Is(err, customerrepo.ErrorDuplicateCustomer) {
			return customerrepo.ErrorDuplicateCustomer
		}
		c.Logger.WithContext(ctx).Error("failed to create customer")
		return fmt.Errorf("creating customer: %w", err)
	}
	data := customersmodels.CommandData{
		Command: "subscribe_by_email",
		Payload: customersmodels.Payload{
			Email:    customer.Email,
			Timezone: customer.Timezone,
		},
	}

	c.Logger.WithContext(ctx).Debug("serializing command")
	serializedData, er := customersmodels.SerializeCommandData(data)
	if er != nil {
		_ = c.CustomerRepo.DeleteCustomerByID(id)
		c.Logger.WithContext(ctx).Error("failed to serialize data")
		return fmt.Errorf("serializing data: %w", er)
	}

	serializedID := fmt.Sprintf("%d", id)

	msg := kafka.Message{
		Key:   []byte(serializedID),
		Value: []byte(serializedData),
		Headers: []kafka.Header{
			{
				Key:   "trace_id",
				Value: []byte(traceID.String()),
			},
		},
	}

	c.Logger.WithContext(ctx).Debug("sending message to broker")
	if err = c.Producer.WriteMessages(ctx, msg); err != nil {
		_ = c.CustomerRepo.DeleteCustomerByID(id)
		c.Logger.WithContext(ctx).Error("failed to send message")
		return fmt.Errorf("sending message: %w", err)
	}

	c.Logger.WithContext(ctx).Debug("message sent")
	return nil
}

func (c *SagaCoordinator) StartReceivingMessages(ctx context.Context) {
	for {
		m, err := c.Consumer.FetchMessage(ctx)

		var traceID string
		for _, h := range m.Headers {
			if h.Key == "trace_id" {
				traceID = string(h.Value)
				break
			}
		}

		ctx = context.WithValue(ctx, logger.TraceIDKey, traceID)
		ctx = context.WithValue(ctx, logger.ServiceNameKey, serviceName)

		if err != nil {
			c.Logger.Error("failed to read message")
			continue
		}

		c.Logger.WithContext(ctx).Info("got message")

		c.Logger.WithContext(ctx).Debug("processing message")
		msg, err := customersmodels.DeserializeReplyData(string(m.Value))
		if err != nil {
			c.Logger.WithContext(ctx).Error("failed to deserialize data")
			continue
		}

		c.Logger.WithContext(ctx).Debug("handling command")
		switch msg.Command {
		case "subscribe_by_email":
			c.Logger.WithContext(ctx).Info("subscribing by email response")
			switch msg.Reply {
			case "success":
				id, er := strconv.Atoi(string(m.Key))
				if er != nil {
					c.Logger.WithContext(ctx).Error("failed to convert id")
					continue
				}
				if err = c.CustomerRepo.UpdateCustomerStatus(id); err != nil {
					c.Logger.WithContext(ctx).Error("failed to update customer status")
					continue
				}

				c.Logger.WithContext(ctx).Info("customer subscribed")
			default:
				c.Logger.WithContext(ctx).Error("unknown reply")
			}

		default:
			c.Logger.WithContext(ctx).Error("unknown command")
		}

		if err = c.Consumer.CommitMessages(ctx, m); err != nil {
			c.Logger.WithContext(ctx).Error("failed to commit message")
		}
	}
}
