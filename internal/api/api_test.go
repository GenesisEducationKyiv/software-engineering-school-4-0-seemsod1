package api

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetUsdToUahRate(t *testing.T) {
	price, err := GetUsdToUahRate()
	assert.NoError(t, err)

	if price == -1 {
		t.Error("expected amount; got empty")
	}
}
