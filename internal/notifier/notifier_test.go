package notifier_test

// type mockEventRepo struct {
//	mock.Mock
// }
//
// func NewMockEventRepo() *mockEventRepo {
//	return &mockEventRepo{}
// }
// func (mr *mockEventRepo) AddToEvents([]emailstreamer.Event) error {
//	args := mr.Called()
//	return args.Error(0)
// }
//
// func (mr *mockEventRepo) ConsumeEvent(event emailstreamer.EventProcessed) error {
//	args := mr.Called(event)
//	return args.Error(0)
// }
//
// func (mr *mockEventRepo) CheckEventProcessed(id int) (bool, error) {
//	args := mr.Called(id)
//	return args.Bool(0), args.Error(1)
//
// }
//
// func TestEmailNotifier_Start_Success(t *testing.T) {
//	mockDB := dbrepo.NewMockDB()
//	eventDB := NewMockEventRepo()
//	mockScheduler := &scheduler.MockScheduler{}
//	provider := &rateapi.CoinbaseProvider{}
//
//	var capturedTask func()
//
//	localTime := time.Now().Hour()
//	timezoneDiff := timezone.GetTimezoneDiff(localTime, notifier.TimeToSend)
//
//	mockScheduler.On("Start").Return(nil)
//	mockScheduler.On("AddEverydayJob", mock.Anything, notifier.MinuteToSend).Return(nil).Run(func(args mock.Arguments) {
//		capturedTask, _ = args.Get(0).(func())
//	})
//
//	os.Setenv("MAILER_HOST", "smtp.gmail.com")
//	os.Setenv("MAILER_PORT", "587")
//	os.Setenv("MAILER_FROM", "abc@mail.com")
//	os.Setenv("MAILER_PASSWORD", "password")
//
//	defer func() {
//		os.Unsetenv("MAILER_HOST")
//		os.Unsetenv("MAILER_PORT")
//		os.Unsetenv("MAILER_FROM")
//		os.Unsetenv("MAILER_PASSWORD")
//	}()
//	logg, _ := logger.NewLogger("test")
//	kafkaProd := emailstreamer.NewKafkaWriter("", "")
//
//	et := notifier.NewEmailNotifier(mockDB, eventDB, mockScheduler, provider, logg, kafkaProd)
//
//	mockDB.On("GetSubscribers", timezoneDiff).Return([]string{}, nil)
//
//	err := et.Start()
//	require.NoError(t, err)
//
//	if capturedTask != nil {
//		capturedTask()
//	} else {
//		t.Fatalf("Scheduled task function was not captured")
//	}
//
//	mockDB.AssertExpectations(t)
//	mockScheduler.AssertExpectations(t)
// }
//
// func TestEmailNotifier_Start_InvalidConfig(t *testing.T) {
//	mockDB := dbrepo.NewMockDB()
//	mockScheduler := &scheduler.MockScheduler{}
//	logg, _ := logger.NewLogger("test")
//
//	os.Setenv("MAILER_HOST", "")
//	os.Setenv("MAILER_PORT", "587")
//	os.Setenv("MAILER_FROM", "")
//
//	defer func() {
//		os.Unsetenv("MAILER_HOST")
//		os.Unsetenv("MAILER_PORT")
//		os.Unsetenv("MAILER_FROM")
//	}()
//
//	et := notifier.NewEmailNotifier(mockDB, mockScheduler, nil, nil, logg)
//	err := et.Start()
//	require.Error(t, err)
// }
//
// func TestEmailNotifier_Start_InvalidConfig_Validation(t *testing.T) {
//	mockDB := dbrepo.NewMockDB()
//	mockScheduler := &scheduler.MockScheduler{}
//	logg, _ := logger.NewLogger("test")
//
//	os.Setenv("MAILER_HOST", "smtp.gmail.com")
//	os.Setenv("MAILER_PORT", "587")
//	os.Setenv("MAILER_FROM", "")
//	os.Setenv("MAILER_PASSWORD", "password")
//
//	defer func() {
//		os.Unsetenv("MAILER_HOST")
//		os.Unsetenv("MAILER_PORT")
//		os.Unsetenv("MAILER_FROM")
//		os.Unsetenv("MAILER_PASSWORD")
//	}()
//	et := notifier.NewEmailNotifier(mockDB, mockScheduler, nil, nil, logg)
//	err := et.Start()
//	require.Error(t, err)
// }
