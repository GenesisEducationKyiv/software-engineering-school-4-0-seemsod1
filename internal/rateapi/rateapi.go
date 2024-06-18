package rateapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"time"
)

type CoinbaseProvider struct{}

func NewCoinbaseProvider() *CoinbaseProvider {
	return &CoinbaseProvider{}
}

const coinbaseURL = "https://api.coinbase.com/v2/prices/%s-%s/buy"

// GetRate returns the current base to target currencies rate
func (cb *CoinbaseProvider) GetRate(base, target string) (float64, error) {
	if !cb.ValidateRateParam(base) || !cb.ValidateRateParam(target) {
		return -1, fmt.Errorf("invalid rate parameters")
	}

	response, err := ProcessGETRequest(fmt.Sprintf(coinbaseURL, base, target))
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

func ProcessGETRequest(url string) ([]byte, error) {
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

func (cb *CoinbaseProvider) ValidateRateParam(code string) bool {
	match, _ := regexp.MatchString("^[A-Z]{3}$", code)
	return match
}
