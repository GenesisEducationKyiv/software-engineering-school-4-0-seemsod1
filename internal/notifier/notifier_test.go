package notifier_test

import (
	"os"
	"testing"
	"time"

	"github.com/seemsod1/api-project/internal/rateapi"

	"github.com/seemsod1/api-project/internal/notifier"
	"github.com/seemsod1/api-project/internal/scheduler"
	"github.com/seemsod1/api-project/internal/storage/dbrepo"
	"github.com/seemsod1/api-project/internal/timezone"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewEmailNotifier(t *testing.T) {
	notification := notifier.NewEmailNotifier(nil, nil, nil, nil)
	require.NotNil(t, notification)
}

func TestEmailNotifier_Start_Success(t *testing.T) {
	mockDB := dbrepo.NewMockDB()
	mockScheduler := &scheduler.MockScheduler{}
	provider := &rateapi.CoinbaseProvider{}

	var capturedTask func()

	localTime := time.Now().Hour()
	timezoneDiff := timezone.GetTimezoneDiff(localTime, notifier.TimeToSend)

	mockScheduler.On("Start").Return(nil)
	mockScheduler.On("AddEverydayJob", mock.Anything, notifier.MinuteToSend).Return(nil).Run(func(args mock.Arguments) {
		capturedTask, _ = args.Get(0).(func())
	})

	os.Setenv("MAILER_HOST", "smtp.gmail.com")
	os.Setenv("MAILER_PORT", "587")
	os.Setenv("MAILER_FROM", "abc@mail.com")
	os.Setenv("MAILER_PASSWORD", "password")

	defer func() {
		os.Unsetenv("MAILER_HOST")
		os.Unsetenv("MAILER_PORT")
		os.Unsetenv("MAILER_FROM")
		os.Unsetenv("MAILER_PASSWORD")
	}()
	cfg, _ := notifier.NewEmailNotifierConfig()
	mailSender, _ := notifier.NewSMTPEmailSender(cfg)
	et := notifier.NewEmailNotifier(mockDB, mockScheduler, provider, mailSender)

	mockDB.On("GetSubscribers", timezoneDiff).Return([]string{}, nil)

	err := et.Start()
	require.NoError(t, err)

	if capturedTask != nil {
		capturedTask()
	} else {
		t.Fatalf("Scheduled task function was not captured")
	}

	mockDB.AssertExpectations(t)
	mockScheduler.AssertExpectations(t)
}

func TestEmailNotifier_Start_InvalidConfig(t *testing.T) {
	mockDB := dbrepo.NewMockDB()
	mockScheduler := &scheduler.MockScheduler{}

	os.Setenv("MAILER_HOST", "")
	os.Setenv("MAILER_PORT", "587")
	os.Setenv("MAILER_FROM", "")

	defer func() {
		os.Unsetenv("MAILER_HOST")
		os.Unsetenv("MAILER_PORT")
		os.Unsetenv("MAILER_FROM")
	}()

	et := notifier.NewEmailNotifier(mockDB, mockScheduler, nil, nil)
	err := et.Start()
	require.Error(t, err)
}

func TestEmailNotifier_Start_InvalidConfig_Validation(t *testing.T) {
	mockDB := dbrepo.NewMockDB()
	mockScheduler := &scheduler.MockScheduler{}

	os.Setenv("MAILER_HOST", "smtp.gmail.com")
	os.Setenv("MAILER_PORT", "587")
	os.Setenv("MAILER_FROM", "")
	os.Setenv("MAILER_PASSWORD", "password")

	defer func() {
		os.Unsetenv("MAILER_HOST")
		os.Unsetenv("MAILER_PORT")
		os.Unsetenv("MAILER_FROM")
		os.Unsetenv("MAILER_PASSWORD")
	}()

	et := notifier.NewEmailNotifier(mockDB, mockScheduler, nil, nil)
	err := et.Start()
	require.Error(t, err)
}
