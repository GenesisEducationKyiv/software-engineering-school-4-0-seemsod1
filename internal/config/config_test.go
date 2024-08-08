package config_test

import (
	"os"
	"testing"

	"github.com/seemsod1/api-project/internal/config"
	"github.com/stretchr/testify/require"
)

func TestNewAppConfig_Success(t *testing.T) {
	os.Setenv("APP_MODE", "test")
	defer os.Unsetenv("APP_MODE")

	appConfig, err := config.NewAppConfig()
	require.NoError(t, err)
	require.NotNil(t, appConfig)
	require.Equal(t, "test", appConfig.Mode)
}
