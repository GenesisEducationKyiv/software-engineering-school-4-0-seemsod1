package main

import (
	"fmt"
	"os"

	logger "github.com/seemsod1/api-project/internal/logger"

	"github.com/seemsod1/api-project/internal/rateapi/chain"

	"github.com/seemsod1/api-project/internal/rateapi"

	"github.com/seemsod1/api-project/internal/scheduler"

	"github.com/seemsod1/api-project/internal/storage/dbrepo"

	"github.com/seemsod1/api-project/internal/config"
	"github.com/seemsod1/api-project/internal/driver"
	"github.com/seemsod1/api-project/internal/handlers"
	"github.com/seemsod1/api-project/internal/notifier"
)

// setup sets up the application
func setup(_ *config.AppConfig, logg *logger.Logger) error {
	dr := driver.NewGORMDriver()

	logg.Info("Connecting to database...")
	db, err := dr.ConnectSQL()
	if err != nil {
		logg.Error("Cannot connect to database! Dying...")
		return fmt.Errorf("connecting to database: %w", err)
	}

	logg.Info("Running migrations...")
	if err = db.RunMigrations(); err != nil {
		logg.Error("Cannot run schemas migration! Dying...")
		return fmt.Errorf("running schemas migration: %w", err)
	}

	CoinBaseProvider := rateapi.NewLoggingClient(os.Getenv("COINBASE_SITE"),
		rateapi.NewCoinbaseProvider(os.Getenv("COINBASE_URL")), logg)

	PrivatBankProvider := rateapi.NewLoggingClient(os.Getenv("PRIVATBANK_SITE"),
		rateapi.NewPrivatBankProvider(os.Getenv("PRIVATBANK_URL")), logg)

	NBUProvider := rateapi.NewLoggingClient(os.Getenv("NBU_SITE"),
		rateapi.NewNBUProvider(os.Getenv("NBU_URL")), logg)

	BaseChain := chain.NewBaseChain(CoinBaseProvider)
	SecondChain := chain.NewBaseChain(PrivatBankProvider)
	ThirdChain := chain.NewBaseChain(NBUProvider)

	BaseChain.SetNext(SecondChain)
	SecondChain.SetNext(ThirdChain)

	dbRepository := dbrepo.NewGormRepo(db.DB)

	repo := handlers.NewRepo(dbRepository, BaseChain, logg)
	handlers.NewHandlers(repo)

	sch := scheduler.NewGoCronScheduler()
	cfg, err := notifier.NewEmailNotifierConfig()
	if err != nil {
		logg.Error("Cannot create mail sender config! Dying...")
		return fmt.Errorf("creating mail sender config: %w", err)
	}

	mailSender, err := notifier.NewSMTPEmailSender(cfg)
	if err != nil {
		logg.Error("Cannot create mail sender! Dying...")
		return fmt.Errorf("creating sender: %w", err)
	}

	logg.Info("Starting mail notifier...")
	notification := notifier.NewEmailNotifier(dbRepository, sch, BaseChain, mailSender, logg)
	if err = notification.Start(); err != nil {
		logg.Error("Cannot start mail notifier! Dying...")
		return fmt.Errorf("starting notifier: %w", err)
	}

	return nil
}
