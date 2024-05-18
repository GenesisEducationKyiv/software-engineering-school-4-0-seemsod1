package helpers

import "testing"

func TestNewRateResponse(t *testing.T) {
	price := 100.0
	response := NewRateResponse(price)
	if response == nil {
		t.Error("expected response; got nil")
	}
}
