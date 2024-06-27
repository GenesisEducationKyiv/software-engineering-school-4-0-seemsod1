package config_test

import (
	"os"
	"testing"

	"github.com/seemsod1/api-project/internal/config"
	"github.com/stretchr/testify/require"
)

func TestNewAppConfig_Success(t *testing.T) {
	os.Setenv("APP_ENV", "true")
	defer os.Unsetenv("APP_ENV")

	appConfig, err := config.NewAppConfig()
	require.NoError(t, err)
	require.NotNil(t, appConfig)
	require.True(t, appConfig.Prod)
}

func TestNewAppConfig_Failure(t *testing.T) {
	os.Clearenv()
	appConfig, err := config.NewAppConfig()
	require.Nil(t, appConfig)
	require.Error(t, err)
}
