package chain_test

import (
	"context"
	"testing"

	"github.com/seemsod1/api-project/internal/rateapi/chain"
	"github.com/stretchr/testify/mock"
)

type mockProvider struct {
	mock.Mock
}

func newMockProvider() *mockProvider {
	return &mockProvider{}
}

func (m *mockProvider) GetRate(ctx context.Context, base, target string) (float64, error) {
	args := m.Called(ctx, base, target)
	return args.Get(0).(float64), args.Error(1)
}

func TestBaseChain_GetRate(t *testing.T) {
	testCases := []struct {
		name        string
		mockReturn  float64
		mockErr     error
		expectedErr error
	}{
		{"ProviderError", 0.0, chain.ErrNoRateProviders, chain.ErrNoRateProviders},
		{"NoError", 27.6, nil, nil},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockProvider := newMockProvider()
			provider := chain.NewBaseChain(mockProvider)

			mockProvider.On("GetRate", context.Background(), "USD", "UAH").Return(tc.mockReturn, tc.mockErr)

			rate, err := provider.GetRate(context.Background(), "USD", "UAH")
			if tc.expectedErr != nil {
				if err == nil || err.Error() != tc.expectedErr.Error() {
					t.Errorf("Expected error %v, got %v", tc.expectedErr, err)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if rate != tc.mockReturn {
					t.Errorf("Expected rate %v, got %v", tc.mockReturn, rate)
				}
			}

			mockProvider.AssertExpectations(t)
		})
	}
}

func TestBaseChain_GetRate_WithNext(t *testing.T) {
	testCases := []struct {
		name            string
		mockReturnFirst float64
		mockErrFirst    error
		mockReturnNext  float64
		mockErrNext     error
		expectedRate    float64
		expectedErr     error
	}{
		{"FirstProviderError", 0.0, chain.ErrNoRateProviders, 27.6, nil, 27.6, nil},
		{"NextProviderError", 0.0, chain.ErrNoRateProviders, 0.0, chain.ErrNoRateProviders, -1, chain.ErrNoRateProviders},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockProvider := newMockProvider()
			provider := chain.NewBaseChain(mockProvider)

			mockProvider.On("GetRate", context.Background(), "USD", "UAH").Return(tc.mockReturnFirst, tc.mockErrFirst)

			mockProvider2 := newMockProvider()
			provider2 := chain.NewBaseChain(mockProvider2)
			provider.SetNext(provider2)

			mockProvider2.On("GetRate", context.Background(), "USD", "UAH").Return(tc.mockReturnNext, tc.mockErrNext)

			rate, err := provider.GetRate(context.Background(), "USD", "UAH")
			if tc.expectedErr != nil {
				if err == nil || err.Error() != tc.expectedErr.Error() {
					t.Errorf("Expected error %v, got %v", tc.expectedErr, err)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if rate != tc.expectedRate {
					t.Errorf("Expected rate %v, got %v", tc.expectedRate, rate)
				}
			}

			// Ensure both mockProvider and mockProvider2 expectations are met
			mockProvider.AssertExpectations(t)
			mockProvider2.AssertExpectations(t)
		})
	}
}
