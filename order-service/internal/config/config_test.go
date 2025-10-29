package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	os.Setenv("APP_ENVIRONMENT", "production")
	os.Setenv("APP_HTTP_SERVER_ADDR", "localhost:8000")
	config := LoadConfig()
	require.NotNil(t, config)
	require.Equal(t, "production", config.Environment)
	require.Equal(t, "localhost:8000", config.HTTPServer.GetAddr())
	os.Setenv("APP_ENVIRONMENT", "staging")
	LoadConfig()
	require.Equal(t, "production", config.Environment)
}
