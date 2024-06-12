package handlers

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/seemsod1/api-project/internal/api"
	"github.com/seemsod1/api-project/internal/helpers"
)

// Rate returns the current USD to UAH rate
func (m *Repository) Rate(w http.ResponseWriter, r *http.Request) {
	price, err := api.GetUsdToUahRate()
	if err != nil {
		http.Error(w, "Failed to get rate", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, helpers.NewRateResponse(price))
}
