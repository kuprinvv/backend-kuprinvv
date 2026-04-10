package auth

import (
	"context"
	"test-backend-1-kuprinvv/internal/model"
	"test-backend-1-kuprinvv/internal/service"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (s *serv) Register(ctx context.Context, email, password, role string) (model.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return model.User{}, err
	}

	user := model.User{
		ID:       uuid.New(),
		Email:    email,
		Password: string(hash),
		Role:     role,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		if service.IsUniqueViolation(err) {
			return model.User{}, model.ErrUserAlreadyExists
		}
		return model.User{}, err
	}

	return user, nil
}
