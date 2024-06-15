package rateapi_test

import (
	"testing"

	"github.com/seemsod1/api-project/internal/rateapi"

	"github.com/stretchr/testify/assert"
)

func TestGetUsdToUahRate(t *testing.T) {
	provider := rateapi.NewCoinbaseProvider()

	price, err := provider.GetRate("USD", "UAH")
	assert.NoError(t, err)

	if price == -1 {
		t.Error("expected amount; got empty")
	}
}
