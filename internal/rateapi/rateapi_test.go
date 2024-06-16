package rateapi_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/seemsod1/api-project/internal/rateapi"

	"github.com/stretchr/testify/assert"
)

func TestCoinbaseProvider_GetRate(t *testing.T) {
	provider := rateapi.NewCoinbaseProvider()

	price, err := provider.GetRate("USD", "UAH")
	require.NoError(t, err)
	require.NotEqual(t, -1, price)
}

func TestCoinbaseProvider_GetRate_InvalidParams(t *testing.T) {
	provider := rateapi.NewCoinbaseProvider()

	price, err := provider.GetRate("USD", "abc")
	require.Error(t, err)
	require.NotEqual(t, -1, price)
}

func TestProcessGETRequest_Success(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{"rate": 27.5}`)
	}))
	defer mockServer.Close()

	base := "USD"
	target := "UAH"

	response, err := rateapi.ProcessGETRequest(fmt.Sprintf("%s/%s/%s", mockServer.URL, base, target))

	assert.NoError(t, err)

	assert.NotNil(t, response)

	expectedResponse := `{"rate": 27.5}`
	require.Contains(t, string(response), expectedResponse)
}

func TestCoinbaseProvider_ValidateRateParam(t *testing.T) {
	tests := []struct {
		name    string
		code    string
		wantErr bool
	}{
		{"valid param", "USD", false},
		{"invalid param - empty", "", true},
		{"invalid params - more than 3 characters", "USDT", true},
		{"invalid params - numbers", "123", true},
		{"invalid params - lowercase", "usd", true},
	}
	provider := rateapi.NewCoinbaseProvider()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := provider.ValidateRateParam(tt.code)
			if tt.wantErr {
				require.False(t, result)
			} else {
				require.True(t, result)
			}
		})
	}
}
