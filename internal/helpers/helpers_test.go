package helpers_test

import (
	"testing"

	"github.com/seemsod1/api-project/internal/helpers"
)

func TestNewRateResponse(t *testing.T) {
	price := 100.0
	response := helpers.NewRateResponse(price)
	if response == nil {
		t.Error("expected response; got nil")
	}
}
