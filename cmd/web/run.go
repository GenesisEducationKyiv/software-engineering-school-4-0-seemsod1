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

	"github.com/seemsod1/api-project/internal/broker"
	emailStreamer "github.com/seemsod1/api-project/internal/email_streamer"
	streamerrepo "github.com/seemsod1/api-project/internal/email_streamer/repository"
	messagesender "github.com/seemsod1/api-project/internal/message_sender"
	messagesenderrepo "github.com/seemsod1/api-project/internal/message_sender/repository"
	notifierrepo "github.com/seemsod1/api-project/internal/notifier/repository"

	"github.com/joho/godotenv"
	"github.com/seemsod1/api-project/internal/config"
	"github.com/seemsod1/api-project/internal/logger"
)

const portNumber = ":8080"

type (
	consumer interface {
		StartReceivingMessages(ctx context.Context)
	}

	producer interface {
		Process(ctx context.Context)
	}
)

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
	kafkaURL := os.Getenv("KAFKA_URL")

	kafReader := broker.NewKafkaConsumer(kafkaURL, "emails", "email_sender_group")

	if err = broker.NewKafkaTopic(kafkaURL, "emails", 1); err != nil {
		return fmt.Errorf("creating kafka topic: %w", err)
	}

	cfg, err := messagesender.NewEmailSenderConfig()
	if err != nil {
		logg.Error("Cannot create mail sender config! Dying...")
		return fmt.Errorf("creating mail sender config: %w", err)
	}

	senderEventRepo, err := messagesenderrepo.NewEventDBRepo(db.DB)
	if err != nil {
		return fmt.Errorf("creating sender event repository: %w", err)
	}

	emailSender, err := messagesender.NewSMTPEmailSender(cfg, kafReader, senderEventRepo, logg)
	if err != nil {
		logg.Error("Cannot create email sender! Dying...")
		return fmt.Errorf("creating email sender: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go eventConsumer(ctx, emailSender)

	kafWriter := broker.NewKafkaProducer(kafkaURL, "emails")

	streamRepository, err := streamerrepo.NewStreamerRepo(db.DB)
	if err != nil {
		return fmt.Errorf("creating streamer repository: %w", err)
	}

	notifierRepository, err := notifierrepo.NewEventDBRepo(db.DB)
	if err != nil {
		return fmt.Errorf("creating notifier repository: %w", err)
	}

	es := emailStreamer.NewEmailStreamer(notifierRepository, streamRepository, kafWriter, logg)
	go eventProducer(ctx, es)

	srv := &http.Server{
		Addr:        portNumber,
		Handler:     routes(),
		ReadTimeout: 30 * time.Second,
	}

	handleShutdown(srv, cancel)

	return nil
}

// handleShutdown handles a graceful shutdown of the application.
func handleShutdown(srv *http.Server, cancelFunc context.CancelFunc) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-stop
		cancelFunc()
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		log.Println("Shutting down server...")
		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("HTTP server shutdown failed: %v", err)
		}
		log.Println("Server has been stopped")
	}()

	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}
}

// eventProducer runs an event dispatcher.
func eventProducer(ctx context.Context, p producer) {
	p.Process(ctx)

	<-ctx.Done()
	log.Println("Shutting down event producer...")
}

// eventProducer runs an event dispatcher.
func eventConsumer(ctx context.Context, c consumer) {
	c.StartReceivingMessages(ctx)

	<-ctx.Done()
	log.Println("Shutting down event consumer...")
}
