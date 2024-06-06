package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/seemsod1/api-project/internal/config"
	"testing"
)

func TestRoutes(t *testing.T) {
	var app config.AppConfig

	mux := routes(&app)

	switch v := mux.(type) {
	case *chi.Mux:
		//do nothing
	default:
		t.Errorf("type is not *chi.Mux, got %v", v)
	}
}
