package main

import (
	"github.com/seemsod1/api-project/internal/config"
	"log"
	"net/http"
	"time"
)

const portNumber = ":8080"

var app config.AppConfig

func main() {
	if err := setup(&app); err != nil {
		log.Fatal(err)
	}

	srv := &http.Server{
		Addr:        portNumber,
		Handler:     routes(),
		ReadTimeout: 30 * time.Second,
	}

	err := srv.ListenAndServe()
	log.Fatal(err)
}
