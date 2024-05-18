package dbrepo

// AddSubscriber adds a new subscriber to the database
func (m *mockDB) AddSubscriber(email string) error {
	args := m.Called(email)
	return args.Error(0)
}

// GetSubscribers returns all subscribers from the database
func (m *mockDB) GetSubscribers() ([]string, error) {
	args := m.Called()
	return args.Get(0).([]string), args.Error(1)
}
