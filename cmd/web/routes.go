package main

import (
	"github.com/VictoriaMetrics/metrics"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/seemsod1/api-project/internal/handlers"
)

// routes sets up the routes for the application
func routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.RequestID)
	//mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)
	mux.Use(EnableCORS)

	mux.Use(Handle)

	mux.Route("/api", func(mux chi.Router) {
		mux.Route("/v1", func(mux chi.Router) {
			mux.Get("/rate", handlers.Repo.Rate)
			mux.Post("/subscribe", handlers.Repo.Subscribe)
			mux.Post("/unsubscribe", handlers.Repo.Unsubscribe)
		})
	})
	mux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		metrics.WritePrometheus(w, true)
	})

	return mux
}
