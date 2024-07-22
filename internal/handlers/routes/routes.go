package routes

import (
	middlewarepkg "github.com/seemsod1/api-project/internal/handlers/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/seemsod1/api-project/internal/handlers"
)

// API sets up the routes for the application
func API(h *handlers.Handlers) http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.RequestID)
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)
	mux.Use(middlewarepkg.EnableCORS)

	mux.Use(middlewarepkg.Handle)

	mux.Route("/api", func(mux chi.Router) {
		mux.Route("/v1", func(mux chi.Router) {
			mux.Get("/rate", h.Rate)
			mux.Post("/subscribe", h.Subscribe)
			mux.Post("/unsubscribe", h.Unsubscribe)
		})
	})

	return mux
}

func Metrics() http.Handler {
	mux := chi.NewRouter()
	mux.HandleFunc("/metrics", handlers.Metrics)
	return mux
}
