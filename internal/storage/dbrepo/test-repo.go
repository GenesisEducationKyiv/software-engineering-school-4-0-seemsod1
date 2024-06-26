package dbrepo

import "github.com/seemsod1/api-project/internal/models"

func (m *MockDB) AddSubscriber(subscriber models.Subscriber) error {
	args := m.Called(subscriber)
	return args.Error(0)
}

func (m *MockDB) GetSubscribers(diff int) ([]string, error) {
	args := m.Called(diff)
	return args.Get(0).([]string), args.Error(1)
}
