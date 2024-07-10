package messagesender_test

// func TestNewSmtpEmailSender_Success(t *testing.T) {
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
//
//	cfg, err := messagesender.NewEmailSenderConfig()
//	require.NoError(t, err)
//
//	require.True(t, cfg.Validate())
//
//	_, err = messagesender.NewSMTPEmailSender(cfg)
//	require.NoError(t, err)
// }
