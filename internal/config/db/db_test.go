package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDBConfig_Success(t *testing.T) {
	t.Setenv(dsnEnvName, "postgres://user:pass@localhost:5432/db")

	cfg, err := NewDBConfig()

	require.NoError(t, err)
	assert.Equal(t, "postgres://user:pass@localhost:5432/db", cfg.DSN())
}

func TestNewDBConfig_MissingEnv(t *testing.T) {
	t.Setenv(dsnEnvName, "")

	_, err := NewDBConfig()

	require.Error(t, err)
}
