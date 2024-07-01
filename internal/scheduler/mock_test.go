package scheduler_test

import (
	"testing"

	"github.com/seemsod1/api-project/internal/scheduler"
	"github.com/stretchr/testify/mock"
)

func TestNewMockScheduler(t *testing.T) {
	ms := scheduler.NewMockScheduler()
	if ms == nil {
		t.Error("expected a non-nil value")
	}
}

func TestMockScheduler_Start(_ *testing.T) {
	ms := scheduler.NewMockScheduler()
	ms.On("Start").Return(nil)
	ms.Start()
}

func TestMockScheduler_AddEverydayJob_Success(t *testing.T) {
	ms := scheduler.NewMockScheduler()
	ms.On("AddEverydayJob", mock.AnythingOfType("func()"), 1).Return(nil)
	err := ms.AddEverydayJob(func() {}, 1)
	if err != nil {
		t.Error("expected no error")
	}
}
