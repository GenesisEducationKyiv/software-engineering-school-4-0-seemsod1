package scheduler_test

import (
	"testing"

	"github.com/seemsod1/api-project/internal/scheduler"
)

func TestNewGoCronScheduler(t *testing.T) {
	gs := scheduler.NewGoCronScheduler()
	if gs == nil {
		t.Error("expected a non-nil value")
	}
}

func TestGoCron_Start(_ *testing.T) {
	gs := scheduler.NewGoCronScheduler()
	gs.Start()
}

func TestGoCron_AddEverydayJob_Success(t *testing.T) {
	gs := scheduler.NewGoCronScheduler()
	err := gs.AddEverydayJob(func() {}, 1)
	if err != nil {
		t.Error("expected no error")
	}
}

func TestGoCron_AddEverydayJob_Failure(t *testing.T) {
	gs := scheduler.NewGoCronScheduler()
	err := gs.AddEverydayJob(func() {}, 1000)
	if err == nil {
		t.Error("expected an error")
	}
}
