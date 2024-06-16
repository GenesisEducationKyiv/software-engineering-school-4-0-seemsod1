package main

import (
	"log"
	"net/http"
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

	err = srv.ListenAndServe()
	log.Fatal(err)
}
