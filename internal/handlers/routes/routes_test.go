package routes_test

import (
	"testing"

	"github.com/seemsod1/api-project/internal/handlers/routes"

	"github.com/go-chi/chi/v5"
)

func TestRoutes(t *testing.T) {
	mux := routes.API(nil)

	switch v := mux.(type) {
	case *chi.Mux:
		// do nothing
	default:
		t.Errorf("type is not *chi.Mux, got %T", v)
	}
}
