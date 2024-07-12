package main

import (
	"fmt"
	"os"

	customer "github.com/seemsod1/api-project/internal/customer"
	customerrepo "github.com/seemsod1/api-project/internal/customer/repository"
	subscriberrepo "github.com/seemsod1/api-project/internal/subscriber/repository"
	"github.com/seemsod1/api-project/pkg/kafkautil"

	"github.com/seemsod1/api-project/internal/handlers"
	"github.com/seemsod1/api-project/internal/notifier"
	notifierrepo "github.com/seemsod1/api-project/internal/notifier/repository"
	"github.com/seemsod1/api-project/internal/rateapi"
	"github.com/seemsod1/api-project/internal/rateapi/chain"
	"github.com/seemsod1/api-project/internal/scheduler"
	"github.com/seemsod1/api-project/pkg/logger"

	"github.com/seemsod1/api-project/internal/config"
	"github.com/seemsod1/api-project/internal/driver"
)

type services struct {
	Driver         *driver.GORMDriver
	Customer       *customer.Service
	SubscriberRepo *subscriberrepo.SubscriberDBRepo
}

// setup sets up the application
func setup(_ *config.AppConfig, logg *logger.Logger) (*services, error) {
	dr := driver.NewGORMDriver()

	logg.Info("Connecting to database...")
	db, err := dr.ConnectSQL()
	if err != nil {
		logg.Error("Cannot connect to database! Dying...")
		return nil, fmt.Errorf("connecting to database: %w", err)
	}

	logg.Info("Setting up rate fetchers chain...")
	fetcher := setupRateFetchersChain(logg)

	subsRepo, err := subscriberrepo.NewSubscriberDBRepo(db.DB)
	if err != nil {
		return nil, fmt.Errorf("creating subscriber repository: %w", err)
	}
	eventRepository, err := notifierrepo.NewEventDBRepo(db.DB)
	if err != nil {
		return nil, fmt.Errorf("creating streamer repository: %w", err)
	}

	sch := scheduler.NewGoCronScheduler()

	logg.Info("Starting mail notifier...")
	notificator := notifier.NewEmailNotifier(subsRepo, eventRepository, sch, fetcher, logg)
	if err = notificator.Start(); err != nil {
		logg.Error("Cannot start mail notifier! Dying...")
		return nil, fmt.Errorf("starting mail notifier: %w", err)
	}
	custRepo, err := customerrepo.NewCustomerRepo(db.DB)
	if err != nil {
		return nil, fmt.Errorf("creating customer repository: %w", err)
	}

	customerKafkaWriter := kafkautil.NewKafkaProducer(os.Getenv("KAFKA_URL"), "subscription")
	customerKafkaReader := kafkautil.NewKafkaConsumer(os.Getenv("KAFKA_URL"), "subscription_responses", "customer_group")

	coordinator := customer.NewSagaCoordinator(custRepo, customerKafkaWriter, customerKafkaReader, logg)

	cust := customer.NewService(custRepo, coordinator)

	repo := handlers.NewRepo(cust.SagaCoordinator, subsRepo, fetcher, logg)
	handlers.NewHandlers(repo)

	return &services{
		Driver:         db,
		Customer:       cust,
		SubscriberRepo: subsRepo,
	}, nil
}

func setupRateFetchersChain(logg *logger.Logger) *chain.Node {
	CoinBaseProvider := rateapi.NewLoggingClient(os.Getenv("COINBASE_SITE"),
		rateapi.NewCoinbaseProvider(os.Getenv("COINBASE_URL")), logg)

	PrivatBankProvider := rateapi.NewLoggingClient(os.Getenv("PRIVATBANK_SITE"),
		rateapi.NewPrivatBankProvider(os.Getenv("PRIVATBANK_URL")), logg)

	NBUProvider := rateapi.NewLoggingClient(os.Getenv("NBU_SITE"),
		rateapi.NewNBUProvider(os.Getenv("NBU_URL")), logg)

	BaseChain := chain.NewNode(CoinBaseProvider)
	SecondChain := chain.NewNode(PrivatBankProvider)
	ThirdChain := chain.NewNode(NBUProvider)

	BaseChain.SetNext(SecondChain)
	SecondChain.SetNext(ThirdChain)

	return BaseChain
}
