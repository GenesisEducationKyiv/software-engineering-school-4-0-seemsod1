package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/seemsod1/api-project/internal/config"
	"github.com/seemsod1/api-project/internal/handlers"
	"net/http"
)

func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	mux.HandleFunc("/rate", handlers.Repo.Rate)

	return mux
}
