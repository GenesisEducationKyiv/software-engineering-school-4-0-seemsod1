package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/render"
	"github.com/seemsod1/api-project/internal/helpers"
	"io"
	"net/http"
)

const url = "https://api.coinbase.com/v2/prices/usd-uah/buy"

func (m *Repository) Rate(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, "Invalid status value", http.StatusBadRequest)
		return
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Invalid status value", http.StatusBadRequest)
		return
	}

	type ApiResponse struct {
		Data struct {
			Amount   string `json:"amount"`
			Base     string `json:"base"`
			Currency string `json:"currency"`
		} `json:"data"`
	}
	var apiResponse ApiResponse
	if err = json.Unmarshal(respBody, &apiResponse); err != nil {
		http.Error(w, "Invalid status value", http.StatusBadRequest)
		return
	}

	var price float64
	if _, err = fmt.Sscanf(apiResponse.Data.Amount, "%f", &price); err != nil {
		http.Error(w, "Invalid status value", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, helpers.NewRateResponse(price))
}
