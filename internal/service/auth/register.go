package auth

import (
	"context"
	"test-backend-1-kuprinvv/internal/model"
	"test-backend-1-kuprinvv/internal/service"

	"github.com/google/uuid"
)

func (s *serv) DummyLogin(ctx context.Context, userID uuid.UUID, role string) (string, error) {
	err := s.userRepo.Create(ctx, model.User{
		ID:       userID,
		Email:    role + "@dummy",
		Password: "dummy",
		Role:     role,
	})
	if err != nil && !service.IsUniqueViolation(err) {
		return "", err
	}

	return s.generateToken(userID, role)
}
