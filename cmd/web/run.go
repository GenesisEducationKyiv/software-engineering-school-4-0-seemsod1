package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/seemsod1/api-project/internal/handlers/routes"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/seemsod1/api-project/internal/subscriber"

	"github.com/seemsod1/api-project/pkg/kafkautil"

	messagesender "github.com/seemsod1/api-project/internal/message_sender"
	messagesenderrepo "github.com/seemsod1/api-project/internal/message_sender/repository"
	notifierrepo "github.com/seemsod1/api-project/internal/notifier/repository"
	emailStreamer "github.com/seemsod1/api-project/pkg/email_streamer"

	"github.com/joho/godotenv"
	"github.com/seemsod1/api-project/internal/config"
	"github.com/seemsod1/api-project/pkg/logger"
)

const (
	apiPortNumber     = ":8080"
	metricsPortNumber = ":8081"
)

type (
	consumer interface {
		StartReceivingMessages(ctx context.Context)
	}

	producer interface {
		Process(ctx context.Context)
	}
)

func run() error {
	_ = godotenv.Load()
	app, err := config.NewAppConfig()
	if err != nil {
		return fmt.Errorf("creating app config: %w", err)
	}
	applogg, err := logger.NewLogger(app.Mode)
	if err != nil {
		return fmt.Errorf("creating logger: %w", err)
	}

	serv, err := setup(app, applogg)
	if err != nil {
		return fmt.Errorf("setting up application: %w", err)
	}
	kafkaURL := os.Getenv("KAFKA_URL")

	kafReader := kafkautil.NewKafkaConsumer(kafkaURL, "emails", "email_sender_group")

	if err = kafkautil.NewKafkaTopic(kafkaURL, "emails", 1); err != nil {
		return fmt.Errorf("creating kafka topic: %w", err)
	}

	if err = kafkautil.NewKafkaTopic(kafkaURL, "subscription", 1); err != nil {
		return fmt.Errorf("creating kafka topic: %w", err)
	}

	if err = kafkautil.NewKafkaTopic(kafkaURL, "subscription_responses", 1); err != nil {
		return fmt.Errorf("creating kafka topic: %w", err)
	}

	cfg, err := messagesender.NewEmailSenderConfig()
	if err != nil {
		applogg.Error("Cannot create mail sender config! Dying...")
		return fmt.Errorf("creating mail sender config: %w", err)
	}

	senderEventRepo, err := messagesenderrepo.NewEventDBRepo(serv.Driver.DB)
	if err != nil {
		return fmt.Errorf("creating sender event repository: %w", err)
	}

	emailSender, err := messagesender.NewSMTPEmailSender(cfg, kafReader, senderEventRepo, applogg)
	if err != nil {
		applogg.Error("Cannot create email sender! Dying...")
		return fmt.Errorf("creating email sender: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go eventConsumer(ctx, emailSender, applogg)

	kafWriter := kafkautil.NewKafkaProducer(kafkaURL, "emails")

	streamRepository, err := notifierrepo.NewStreamerRepo(serv.Driver.DB)
	if err != nil {
		return fmt.Errorf("creating streamer repository: %w", err)
	}

	notifierRepository, err := notifierrepo.NewEventDBRepo(serv.Driver.DB)
	if err != nil {
		return fmt.Errorf("creating notifier repository: %w", err)
	}

	es := emailStreamer.NewEmailStreamer(notifierRepository, streamRepository, kafWriter, applogg)
	go eventProducer(ctx, es, applogg)

	subscriberKafkaWriter := kafkautil.NewKafkaProducer(kafkaURL, "subscription_responses")
	subscriberKafkaReader := kafkautil.NewKafkaConsumer(kafkaURL, "subscription", "subscriber_group")
	subs := subscriber.NewService(serv.SubscriberRepo, subscriberKafkaWriter, subscriberKafkaReader, applogg)
	go eventConsumer(ctx, subs, applogg)

	go eventConsumer(ctx, serv.Customer.SagaCoordinator, applogg)

	apiSrv := &http.Server{
		Addr:        apiPortNumber,
		Handler:     routes.API(serv.Handlers),
		ReadTimeout: 30 * time.Second,
	}

	metricsSrv := &http.Server{
		Addr:        metricsPortNumber,
		Handler:     routes.Metrics(),
		ReadTimeout: 30 * time.Second,
	}

	handleShutdown(apiSrv, metricsSrv, cancel, applogg)

	return nil
}

// handleShutdown handles a graceful shutdown of the application.
func handleShutdown(apiSrv *http.Server, metricsSrv *http.Server, cancelFunc context.CancelFunc, l *logger.Logger) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-stop
		cancelFunc()
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		l.Info("Shutting down servers...")

		if err := apiSrv.Shutdown(ctx); err != nil {
			l.Error("API server shutdown failed", zap.Error(err))
		}

		if err := metricsSrv.Shutdown(ctx); err != nil {
			l.Error("Metrics server shutdown failed", zap.Error(err))
		}

		l.Info("Servers has been stopped")
	}()

	l.Info("Starting servers...")
	l.Info("API server is running on", zap.String("port", apiSrv.Addr))
	go func() {
		if err := apiSrv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			l.Fatal("failed to shutdown API server", zap.Error(err))
		}
	}()

	l.Info("Metrics server is running on", zap.String("port", metricsSrv.Addr))
	if err := metricsSrv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		l.Fatal("Metrics server ListenAndServe", zap.Error(err))
	}
}

// eventProducer runs an event dispatcher.
func eventProducer(ctx context.Context, p producer, l *logger.Logger) {
	p.Process(ctx)

	<-ctx.Done()
	l.Info("Shutting down event producer...")
}

// eventProducer runs an event dispatcher.
func eventConsumer(ctx context.Context, c consumer, l *logger.Logger) {
	c.StartReceivingMessages(ctx)

	<-ctx.Done()
	l.Info("Shutting down event consumer...")
}
