package auth

import (
	"test-backend-1-kuprinvv/internal/config"
	"test-backend-1-kuprinvv/internal/repository"
	"test-backend-1-kuprinvv/internal/service"
)

var _ service.AuthService = (*serv)(nil)

type serv struct {
	jwtConf  config.JWTConfig
	userRepo repository.UserRepository
}

func NewAuthService(jwtConfig config.JWTConfig, userRepo repository.UserRepository) *serv {
	return &serv{jwtConf: jwtConfig, userRepo: userRepo}
}
