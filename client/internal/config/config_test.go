package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	envFileContent = []byte(
		`SERVER_ADDRESS=localhost:1234`,
	)
	envFilePath = ".env.test"

	expectedAppConfig = AppConfig{
		ServerAddress: "localhost:1234",
	}
)

func TestNewConfigFromEnvFile(t *testing.T) {
	createFakeEnvFile(t)
	defer deleteFakeEnvFile(t)

	config, err := NewConfig(envFilePath)
	require.NoError(t, err)

	require.Equal(t, expectedAppConfig, config.AppConfig)
}

func createFakeEnvFile(t *testing.T) {
	t.Helper()
	if err := os.WriteFile(envFilePath, envFileContent, 0600); err != nil { //nolint:gofumpt
		t.Fatal(err)
	}
}

func deleteFakeEnvFile(t *testing.T) {
	t.Helper()
	if err := os.Remove(envFilePath); err != nil {
		t.Fatal(err)
	}
}
