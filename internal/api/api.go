package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const url = "https://api.coinbase.com/v2/prices/usd-uah/buy"

// GetUsdToUahRate returns the current USD to UAH rate
func GetUsdToUahRate() (float64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client := &http.Client{}
	var err error

	req, err := http.NewRequestWithContext(ctx, "GET", url, http.NoBody)
	if err != nil {
		return -1, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return -1, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return -1, err
	}

	fmt.Println("Response Body:", string(respBody))

	type APIResponse struct {
		Data struct {
			Amount   string `json:"amount"`
			Base     string `json:"base"`
			Currency string `json:"currency"`
		} `json:"data"`
	}
	var apiResponse APIResponse

	if marshalErr := json.Unmarshal(respBody, &apiResponse); marshalErr != nil {
		return -1, marshalErr
	}

	var price float64
	if _, err = fmt.Sscanf(apiResponse.Data.Amount, "%f", &price); err != nil {
		return -1, err
	}

	return price, nil
}
