package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig_Success(t *testing.T) {
	// temporary config file
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "test_config.yaml")

	configContent := `service:
  host: localhost
  port: 8080
  file: "/output/task-db.json"
  interval: 3
`

	err := os.WriteFile(configFile, []byte(configContent), 0644)
	require.NoError(t, err)

	cfg, err := LoadConfig(configFile)

	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "localhost", cfg.Service.Host)
	assert.Equal(t, 8080, cfg.Service.Port)
	assert.Equal(t, "/output/task-db.json", cfg.Service.File)
	assert.Equal(t, 3, cfg.Service.Interval)
}
