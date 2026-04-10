package jwt

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewJwtConfig_Success(t *testing.T) {
	t.Setenv(jwtEnvName, "supersecretkey")

	cfg, err := NewJwtConfig()

	require.NoError(t, err)
	assert.Equal(t, "supersecretkey", cfg.Token())
}

func TestNewJwtConfig_MissingEnv(t *testing.T) {
	t.Setenv(jwtEnvName, "")

	_, err := NewJwtConfig()

	require.Error(t, err)
}
