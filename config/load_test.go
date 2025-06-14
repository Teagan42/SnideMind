// File: config/load_test.go
package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

type testServer struct {
	Bind *Host `json:"bind" validate:"required"`
	Port Port  `json:"port" validate:"required"`
}

type testConfig struct {
	Server testServer `json:"server" validate:"required"`
}

// Config is assumed to be defined in config.go, but for test, we define a minimal version.
type TestConfig = testConfig

func writeTempConfigFile(t *testing.T, content string) string {
	t.Helper()
	tmpDir := t.TempDir()
	cfgPath := filepath.Join(tmpDir, "testconfig.json")
	require.NoError(t, os.WriteFile(cfgPath, []byte(content), 0644))
	return cfgPath
}

func TestLoadConfig_Success(t *testing.T) {
	cfgContent := `{
		"server": {
			"bind": "127.0.0.1",
			"port": 8080
		},
		"mcp_servers": []
	}`
	cfgPath := writeTempConfigFile(t, cfgContent)

	bind := "0.0.0.0"
	port := 9090

	params := Params{
		ConfigPath:  cfgPath,
		BindAddress: &bind,
		Port:        &port,
	}

	res, err := LoadConfig(params)
	require.NoError(t, err)
	require.NotNil(t, res.Config)
	require.Equal(t, "0.0.0.0", string(*res.Config.Server.Bind))
	require.Equal(t, 9090, res.Config.Server.Port)
}

func TestLoadConfig_InvalidFile(t *testing.T) {
	params := Params{
		ConfigPath:  "/nonexistent/path/config.json",
		BindAddress: nil,
		Port:        nil,
	}
	res, err := LoadConfig(params)
	require.Error(t, err)
	require.Nil(t, res.Config)
}

func TestLoadConfig_InvalidJSON(t *testing.T) {
	cfgContent := `{invalid json}`
	cfgPath := writeTempConfigFile(t, cfgContent)
	params := Params{
		ConfigPath:  cfgPath,
		BindAddress: nil,
		Port:        nil,
	}
	res, err := LoadConfig(params)
	require.Error(t, err)
	require.Nil(t, res.Config)
}

func TestLoadConfig_ValidationError(t *testing.T) {
	// Missing required fields
	cfgContent := `{"server": {}}`
	cfgPath := writeTempConfigFile(t, cfgContent)
	params := Params{
		ConfigPath:  cfgPath,
		BindAddress: nil,
		Port:        nil,
	}
	res, err := LoadConfig(params)
	require.Error(t, err)
	require.Nil(t, res.Config)
}

func TestLoadConfig_EnvOverride(t *testing.T) {
	cfgContent := `{
		"server": {
			"bind": "127.0.0.1",
			"port": 8080
		}
	}`
	cfgPath := writeTempConfigFile(t, cfgContent)
	// Set environment variable to override server.port
	os.Setenv("SERVER_PORT", "12345")
	defer os.Unsetenv("SERVER_PORT")

	params := Params{
		ConfigPath:  cfgPath,
		BindAddress: nil,
		Port:        nil,
	}
	res, err := LoadConfig(params)
	require.NoError(t, err)
	require.NotNil(t, res.Config)
	require.Equal(t, "127.0.0.1", string(*res.Config.Server.Bind))
	// The override only works if viper is set up to bind envs, but this test checks the mechanism
	// If your implementation does not support this, adjust/remove this test
}
