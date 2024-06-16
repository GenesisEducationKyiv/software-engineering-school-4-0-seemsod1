package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"github.com/seemsod1/api-project/internal/config"
)

const portNumber = ":8080"

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	app, err := config.NewAppConfig()
	if err != nil {
		log.Fatal(err)
	}

	if err = setup(app); err != nil {
		log.Fatal(err)
	}

	srv := &http.Server{
		Addr:        portNumber,
		Handler:     routes(),
		ReadTimeout: 30 * time.Second,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("Server is running on port %s\n", portNumber)
		if err = srv.ListenAndServe(); err != nil && !errors.Is(http.ErrServerClosed, err) {
			log.Fatalf("Error starting server: %v\n", err)
		}
	}()

	<-stop

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err = srv.Shutdown(ctx); err != nil {
		log.Printf("Error shutting down server: %v\n", err)
	}

	log.Println("Server shutdown complete.")
}
