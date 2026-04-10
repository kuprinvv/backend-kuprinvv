package auth

import (
	"test-backend-1-kuprinvv/internal/service"
)

type Handler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *Handler {
	return &Handler{authService: authService}
}
