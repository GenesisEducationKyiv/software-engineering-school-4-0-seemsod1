package scheduler

import "github.com/stretchr/testify/mock"

type MockScheduler struct {
	mock.Mock
}

func NewMockScheduler() *MockScheduler {
	return &MockScheduler{}
}

func (ms *MockScheduler) Start() {
	ms.Called()
}

func (ms *MockScheduler) AddEverydayJob(task func(), minute int) error {
	args := ms.Called(task, minute)
	return args.Error(0)
}
