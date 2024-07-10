package handlers_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/seemsod1/api-project/pkg/logger"

	"github.com/seemsod1/api-project/internal/handlers"
	"github.com/stretchr/testify/assert"
)

func TestRate_ProviderError(t *testing.T) {
	mockProvider := newMockProvider()
	logg, _ := logger.NewLogger("test")
	repo := handlers.Repository{RateService: mockProvider, Logger: logg}

	handler := http.HandlerFunc(repo.Rate)

	mockProvider.On("GetRate", context.Background(), "USD", "UAH").Return(0.0, errors.New("provider error"))

	req := httptest.NewRequest("GET", "/rate", http.NoBody)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	mockProvider.AssertExpectations(t)
}

func TestRate_Success(t *testing.T) {
	provider := newMockProvider()
	provider.On("GetRate", context.Background(), "USD", "UAH").Return(27.6, nil)

	repo := handlers.Repository{RateService: provider}

	handler := http.HandlerFunc(repo.Rate)

	req := httptest.NewRequest("GET", "/rate", http.NoBody)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), `"price":`)
}
