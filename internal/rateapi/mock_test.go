package rateapi_test

import (
	"testing"

	"github.com/seemsod1/api-project/internal/rateapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMockProvider_GetRate(t *testing.T) {
	provider := rateapi.NewMockProvider()

	provider.ExpectedCalls = nil
	provider.On("GetRate", "USD", "UAH").Return(1.0, nil)

	price, err := provider.GetRate("USD", "UAH")
	require.NoError(t, err)
	require.Equal(t, 1.0, price)

	provider.AssertExpectations(t)
}

func TestNewMockProvider(t *testing.T) {
	provider := rateapi.NewMockProvider()

	assert.NotNil(t, provider)
}
