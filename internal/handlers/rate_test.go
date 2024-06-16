package handlers_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/seemsod1/api-project/internal/handlers"
	"github.com/seemsod1/api-project/internal/rateapi"
	"github.com/stretchr/testify/assert"
)

func TestRate_ProviderError(t *testing.T) {
	mockProvider := rateapi.NewMockProvider()
	repo := handlers.Repository{Provider: mockProvider}

	handler := http.HandlerFunc(repo.Rate)

	mockProvider.On("GetRate", "USD", "UAH").Return(0.0, errors.New("provider error"))

	req := httptest.NewRequest("GET", "/rate", http.NoBody)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	mockProvider.AssertExpectations(t)
}

func TestRate_Success(t *testing.T) {
	provider := rateapi.NewCoinbaseProvider()
	repo := handlers.Repository{Provider: provider}

	handler := http.HandlerFunc(repo.Rate)

	req := httptest.NewRequest("GET", "/rate", http.NoBody)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), `"price":`)
}
