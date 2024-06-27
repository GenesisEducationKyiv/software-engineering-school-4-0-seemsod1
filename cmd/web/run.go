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
	defer logg.Sync()

	if err = setup(app, logg); err != nil {
		return fmt.Errorf("setting up application: %w", err)
	}

	srv := &http.Server{
		Addr:        portNumber,
		Handler:     routes(),
		ReadTimeout: 30 * time.Second,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		logg.Infof("Server is running on port %s\n", portNumber)
		if err = srv.ListenAndServe(); err != nil && !errors.Is(http.ErrServerClosed, err) {
			logg.Fatalf("Error starting server: %v\n", err)
		}
	}()

	<-stop

	logg.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err = srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("shutting down server: %w", err)
	}

	logg.Info("Server shutdown complete.")
	return nil
}
