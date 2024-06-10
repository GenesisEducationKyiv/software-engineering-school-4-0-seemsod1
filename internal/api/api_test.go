package api_test

import (
	"testing"

	"github.com/seemsod1/api-project/internal/api"

	"github.com/stretchr/testify/assert"
)

func TestGetUsdToUahRate(t *testing.T) {
	price, err := api.GetUsdToUahRate()
	assert.NoError(t, err)

	if price == -1 {
		t.Error("expected amount; got empty")
	}
}
