package handlers

import (
	"net/http"

	"github.com/seemsod1/api-project/internal/rateapi"

	"github.com/go-chi/render"
)

// Rate returns the current USD to UAH rate
func (m *Repository) Rate(w http.ResponseWriter, r *http.Request) {
	provider := rateapi.NewCoinbaseProvider()

	price, err := provider.GetRate("USD", "UAH")
	if err != nil {
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
