package rateapi

import "github.com/stretchr/testify/mock"

type MockProvider struct {
	mock.Mock
}

func NewMockProvider() *MockProvider {
	return &MockProvider{}
}

func (mp *MockProvider) GetRate(base, target string) (float64, error) {
	args := mp.Called(base, target)
	return args.Get(0).(float64), args.Error(1)
}
