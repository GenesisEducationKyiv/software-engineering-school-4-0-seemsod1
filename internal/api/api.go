package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const url = "https://api.coinbase.com/v2/prices/usd-uah/buy"

// GetUsdToUahRate returns the current USD to UAH rate
func GetUsdToUahRate() (float64, error) {
	resp, err := http.Get(url)
	if err != nil {
		return -1, err
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return -1, err
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
		return -1, err
	}

	var price float64
	if _, err = fmt.Sscanf(apiResponse.Data.Amount, "%f", &price); err != nil {
		return -1, err
	}

	return price, nil
}
