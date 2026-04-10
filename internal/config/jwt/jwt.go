package jwt

import (
	"fmt"
	"os"
	"test-backend-1-kuprinvv/internal/config"
)

var _ config.JWTConfig = (*jwtConfig)(nil)

const (
	jwtEnvName = "JWT"
)

type jwtConfig struct {
	jwt string
}

func NewJwtConfig() (*jwtConfig, error) {
	jwt := os.Getenv(jwtEnvName)
	if len(jwt) == 0 {
		return nil, fmt.Errorf("jwt env var not set")
	}

	return &jwtConfig{jwt: jwt}, nil
}

func (j *jwtConfig) Token() string {
	return j.jwt
}
