package main

import (
	"fmt"
	"os"

	"github.com/seemsod1/api-project/internal/handlers"
	"github.com/seemsod1/api-project/internal/notifier"
	notifierrepo "github.com/seemsod1/api-project/internal/notifier/repository"
	"github.com/seemsod1/api-project/internal/rateapi"
	"github.com/seemsod1/api-project/internal/rateapi/chain"
	"github.com/seemsod1/api-project/internal/scheduler"
	"github.com/seemsod1/api-project/internal/storage/dbrepo"

	logger "github.com/seemsod1/api-project/internal/logger"

	"github.com/seemsod1/api-project/internal/config"
	"github.com/seemsod1/api-project/internal/driver"
)

// setup sets up the application
func setup(_ *config.AppConfig, logg *logger.Logger) (*driver.GORMDriver, error) {
	dr := driver.NewGORMDriver()

	logg.Info("Connecting to database...")
	db, err := dr.ConnectSQL()
	if err != nil {
		logg.Error("Cannot connect to database! Dying...")
		return nil, fmt.Errorf("connecting to database: %w", err)
	}

	logg.Info("Running migrations...")
	if err = db.RunMigrations(); err != nil {
		logg.Error("Cannot run schemas migration! Dying...")
		return nil, fmt.Errorf("running schemas migration: %w", err)
	}

	logg.Info("Setting up rate fetchers chain...")
	fetcher := setupRateFetchersChain(logg)

	dbRepository := dbrepo.NewGormRepo(db.DB)
	eventRepository, err := notifierrepo.NewEventDBRepo(db.DB)
	if err != nil {
		return nil, fmt.Errorf("creating streamer repository: %w", err)
	}

	sch := scheduler.NewGoCronScheduler()

	logg.Info("Starting mail notifier...")
	notificator := notifier.NewEmailNotifier(dbRepository, eventRepository, sch, fetcher, logg)
	if err = notificator.Start(); err != nil {
		logg.Error("Cannot start mail notifier! Dying...")
		return nil, fmt.Errorf("starting mail notifier: %w", err)
	}

	repo := handlers.NewRepo(dbRepository, fetcher, notificator, logg)
	handlers.NewHandlers(repo)

	return db, nil
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
