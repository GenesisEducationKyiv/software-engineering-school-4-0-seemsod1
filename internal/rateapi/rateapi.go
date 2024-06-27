package rateapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"time"
)

type CoinbaseProvider struct {
	URL string
}

func NewCoinbaseProvider(url string) *CoinbaseProvider {
	return &CoinbaseProvider{
		URL: url,
	}
}

// GetRate returns the current base to target currencies rate
func (cb *CoinbaseProvider) GetRate(ctx context.Context, base, target string) (float64, error) {
	if !cb.ValidateRateParam(base) || !cb.ValidateRateParam(target) {
		return -1, fmt.Errorf("invalid rate parameters")
	}

	response, err := ProcessGETRequest(ctx, fmt.Sprintf(cb.URL, base, target))
	if err != nil {
		return -1, fmt.Errorf("process get request: %w", err)
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
		return -1, fmt.Errorf("unmarshaling response: %w", marshalErr)
	}

	var price float64
	if _, err = fmt.Sscanf(coinbaseAPIResponse.Data.Amount, "%f", &price); err != nil {
		return -1, fmt.Errorf("parsing price: %w", err)
	}

	return price, nil
}

func ProcessGETRequest(ctx context.Context, url string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	client := &http.Client{}
	var err error

	req, err := http.NewRequestWithContext(ctx, "GET", url, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("making request: %w", err)
	}
	defer func() {
		if err = resp.Body.Close(); err != nil {
			log.Printf("error closing response body: %v", err)
		}
	}()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	return respBody, nil
}

func (cb *CoinbaseProvider) ValidateRateParam(code string) bool {
	match, _ := regexp.MatchString("^[A-Z]{3}$", code)
	return match
}

// PrivatBankProvider is a provider for fetching exchange rates from PrivatBank
type PrivatBankProvider struct {
	URL string
}

func NewPrivatBankProvider(url string) *PrivatBankProvider {
	return &PrivatBankProvider{
		URL: url,
	}
}

// GetRate returns the current base to target currencies rate
func (pb *PrivatBankProvider) GetRate(ctx context.Context, base, target string) (float64, error) {
	if !pb.ValidateRateParam(base) || !pb.ValidateRateParam(target) {
		return -1, fmt.Errorf("invalid rate parameters")
	}

	response, err := ProcessGETRequest(ctx, pb.URL)
	if err != nil {
		return -1, fmt.Errorf("process get request: %w", err)
	}

	type PrivatBankAPIResponse struct {
		Ccy      string `json:"ccy"`
		BaseCcy  string `json:"base_ccy"`
		BuyRate  string `json:"buy"`
		SaleRate string `json:"sale"`
	}
	var privatBankAPIResponse []PrivatBankAPIResponse

	if marshalErr := json.Unmarshal(response, &privatBankAPIResponse); marshalErr != nil {
		return -1, fmt.Errorf("unmarshaling response: %w", marshalErr)
	}

	for _, rate := range privatBankAPIResponse {
		if rate.Ccy == base && rate.BaseCcy == target {
			var price float64
			if _, err = fmt.Sscanf(rate.BuyRate, "%f", &price); err != nil {
				return -1, fmt.Errorf("parsing price: %w", err)
			}
			return price, nil
		}
	}

	return -1, fmt.Errorf("get rate: rate not found")
}

func (pb *PrivatBankProvider) ValidateRateParam(code string) bool {
	match, _ := regexp.MatchString("^[A-Z]{3}$", code)
	return match
}

// NBUProvider is a provider for fetching exchange rates from the National Bank of Ukraine
type NBUProvider struct {
	URL string
}

func NewNBUProvider(url string) *NBUProvider {
	return &NBUProvider{
		URL: url,
	}
}

// GetRate returns the current base to target currencies rate
func (nbu *NBUProvider) GetRate(ctx context.Context, base, target string) (float64, error) {
	if !nbu.ValidateRateParam(base) || !nbu.ValidateRateParam(target) {
		return -1, fmt.Errorf("invalid rate parameters")
	}

	response, err := ProcessGETRequest(ctx, fmt.Sprintf(nbu.URL, base))
	if err != nil {
		return -1, fmt.Errorf("process get request: %w", err)
	}

	type NBUAPIResponse struct {
		Rate float64 `json:"rate"`
	}
	var nbuAPIResponse []NBUAPIResponse

	if marshalErr := json.Unmarshal(response, &nbuAPIResponse); marshalErr != nil {
		return -1, fmt.Errorf("unmarshaling response: %w", marshalErr)
	}

	return nbuAPIResponse[0].Rate, nil
}

func (nbu *NBUProvider) ValidateRateParam(code string) bool {
	match, _ := regexp.MatchString("^[A-Z]{3}$", code)
	return match
}
