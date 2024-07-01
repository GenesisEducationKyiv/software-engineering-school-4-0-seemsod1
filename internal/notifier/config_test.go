package notifier_test

import (
	"os"
	"testing"

	"github.com/seemsod1/api-project/internal/notifier"
	"github.com/stretchr/testify/require"
)

func TestEmailNotifierConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		c       *notifier.EmailNotifierConfig
		wantErr bool
	}{
		{
			name:    "valid config",
			c:       &notifier.EmailNotifierConfig{Host: "smtp.gmail.com", Port: "587", From: "example@gmail.com", Password: "password"},
			wantErr: false,
		},
		{
			name:    "invalid config - missing From",
			c:       &notifier.EmailNotifierConfig{Host: "smtp.gmail.com", Port: "587", From: "", Password: "password"},
			wantErr: true,
		},
		{
			name:    "invalid config - missing Host",
			c:       &notifier.EmailNotifierConfig{Host: "", Port: "587", From: "abc", Password: "password"},
			wantErr: true,
		},
		{
			name:    "invalid config - missing Port",
			c:       &notifier.EmailNotifierConfig{Host: "smtp.gmail.com", Port: "", From: "abc", Password: "password"},
			wantErr: true,
		},
		{
			name:    "invalid config - missing Password",
			c:       &notifier.EmailNotifierConfig{Host: "smtp.gmail.com", Port: "587", From: "abc", Password: ""},
			wantErr: true,
		},
		{
			name:    "invalid config - missing all",
			c:       &notifier.EmailNotifierConfig{Host: "", Port: "", From: "", Password: ""},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.c.Validate()
			if tt.wantErr {
				require.False(t, result)
			} else {
				require.True(t, result)
			}
		})
	}
}

func TestNewEmailNotifierConfig_Invalid(t *testing.T) {
	_, err := notifier.NewEmailNotifierConfig()
	require.Error(t, err)
}

func TestNewEmailNotifierConfig_Valid(t *testing.T) {
	_ = os.Setenv("MAILER_HOST", "smtp.gmail.com")
	_ = os.Setenv("MAILER_PORT", "587")
	_ = os.Setenv("MAILER_FROM", "xad@mail.com")
	_ = os.Setenv("MAILER_PASSWORD", "password")

	defer func() {
		_ = os.Unsetenv("MAILER_HOST")
		_ = os.Unsetenv("MAILER_PORT")
		_ = os.Unsetenv("MAILER_FROM")
		_ = os.Unsetenv("MAILER_PASSWORD")
	}()

	_, err := notifier.NewEmailNotifierConfig()
	require.NoError(t, err)
}
