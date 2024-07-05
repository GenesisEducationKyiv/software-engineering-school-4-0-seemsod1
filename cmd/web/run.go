package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	emailStreamer "github.com/seemsod1/api-project/internal/email_streamer"
	streamerrepo "github.com/seemsod1/api-project/internal/email_streamer/repository"
	"github.com/seemsod1/api-project/internal/handlers"
	messagesender "github.com/seemsod1/api-project/internal/message_sender"
	"github.com/seemsod1/api-project/internal/notifier"
	"github.com/seemsod1/api-project/internal/rateapi"
	"github.com/seemsod1/api-project/internal/rateapi/chain"
	"github.com/seemsod1/api-project/internal/scheduler"
	"github.com/seemsod1/api-project/internal/storage/dbrepo"

	"github.com/joho/godotenv"
	"github.com/seemsod1/api-project/internal/config"
	"github.com/seemsod1/api-project/internal/logger"
)

const portNumber = ":8080"

func run() error {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	app, err := config.NewAppConfig()
	if err != nil {
		return fmt.Errorf("creating app config: %w", err)
	}
	logg, err := logger.NewLogger(app.Mode)
	if err != nil {
		return fmt.Errorf("creating logger: %w", err)
	}

	db, err := setup(app, logg)
	if err != nil {
		return fmt.Errorf("setting up application: %w", err)
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
	streamRepository, err := streamerrepo.NewStreamerRepo(db.DB)
	if err != nil {
		return fmt.Errorf("creating streamer repository: %w", err)
	}

	kafkaURL := os.Getenv("KAFKA_URL")

	kafReader := emailStreamer.NewKafkaConsumer(kafkaURL, "emails", "email_sender_group")
	if err = emailStreamer.NewKafkaTopic(kafkaURL, "emails"); err != nil {
		return fmt.Errorf("creating kafka topic: %w", err)
	}

	cfg, err := messagesender.NewEmailSenderConfig()
	if err != nil {
		logg.Error("Cannot create mail sender config! Dying...")
		return fmt.Errorf("creating mail sender config: %w", err)
	}

	emailSender, err := messagesender.NewSMTPEmailSender(cfg, kafReader, streamRepository, logg)
	if err != nil {
		logg.Error("Cannot create email sender! Dying...")
		return fmt.Errorf("creating email sender: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go emailSender.StartReceivingMessages(ctx)

	kafWriter := emailStreamer.NewKafkaWriter(kafkaURL, "emails")

	es := emailStreamer.NewEmailStreamer(streamRepository, kafWriter, logg)
	go es.Process(ctx)

	sch := scheduler.NewGoCronScheduler()

	logg.Info("Starting mail notifier...")

	notification := notifier.NewEmailNotifier(dbRepository, streamRepository, sch, BaseChain, logg, kafWriter)

	repo := handlers.NewRepo(dbRepository, BaseChain, notification, logg)
	handlers.NewHandlers(repo)

	if err = notification.Start(); err != nil {
		logg.Error("Cannot start mail notifier! Dying...")
		return fmt.Errorf("starting notifier: %w", err)
	}

	srv := &http.Server{
		Addr:        portNumber,
		Handler:     routes(),
		ReadTimeout: 30 * time.Second,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		logg.Infof("Server is running on port %s", portNumber)
		if err = srv.ListenAndServe(); err != nil && !errors.Is(http.ErrServerClosed, err) {
			logg.Fatalf("Error starting server: %v", err)
		}
	}()

	<-stop

	logg.Info("Shutting down server...")

	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err = srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("shutting down server: %w", err)
	}

	logg.Info("Server shutdown complete.")
	return nil
}
