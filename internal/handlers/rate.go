package handlers

import (
	"context"
	"net/http"

	"go.uber.org/zap"

	"github.com/go-chi/render"
)

type RateService interface {
	GetRate(ctx context.Context, base, target string) (float64, error)
}

// Rate returns the current USD to UAH rate
func (m *Repository) Rate(w http.ResponseWriter, r *http.Request) {
	price, err := m.RateService.GetRate(r.Context(), "USD", "UAH")
	if err != nil {
		m.Logger.Error("Getting rate", zap.Error(err))
		http.Error(w, "Failed to get rate", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, newRateResponse(price))
}

// newRateResponse is a helper function that creates a new RateResponse struct
func newRateResponse(price float64) interface{} {
	type RateResponse struct {
		Price float64 `json:"price"`
	}
	return RateResponse{Price: price}
}
