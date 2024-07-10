package dbrepo

import "github.com/seemsod1/api-project/internal/models"

func (m *MockDB) AddSubscriber(subscriber models.Subscriber) error {
	args := m.Called(subscriber)
	return args.Error(0)
}

func (m *MockDB) GetSubscribersWithTimezone(diff int) ([]string, error) {
	args := m.Called(diff)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockDB) GetSubscribers() ([]string, error) {
	args := m.Called()
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockDB) RemoveSubscriber(email string) error {
	args := m.Called(email)
	return args.Error(0)
}
