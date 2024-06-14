package rateapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type CoinbaseProvider struct{}

func NewCoinbaseProvider() *CoinbaseProvider {
	return &CoinbaseProvider{}
}

const coinbaseURL = "https://api.coinbase.com/v2/prices/usd-uah/buy"

// GetUsdToUahRate returns the current USD to UAH rate
func (cb *CoinbaseProvider) GetUsdToUahRate() (float64, error) {
	response, err := processGETRequest(coinbaseURL)
	if err != nil {
		return -1, err
	}

	type CoinbaseAPIResponse struct {
		Data struct {
			Amount   string `json:"amount"`
			Base     string `json:"base"`
			Currency string `json:"currency"`
		} `json:"data"`
	}
	var coinbaseAPIResponse CoinbaseAPIResponse

	if marshalErr := json.Unmarshal(response, &coinbaseAPIResponse); marshalErr != nil {
		return -1, marshalErr
	}

	var price float64
	if _, err = fmt.Sscanf(coinbaseAPIResponse.Data.Amount, "%f", &price); err != nil {
		return -1, err
	}

	return price, nil
}

func processGETRequest(url string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client := &http.Client{}
	var err error

	req, err := http.NewRequestWithContext(ctx, "GET", url, http.NoBody)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return respBody, nil
}
