package converter

import (
	"test-backend-1-kuprinvv/internal/handler/auth/dto"
	"test-backend-1-kuprinvv/internal/model"
)

func ServiceToRegisterResponse(user model.User) dto.RegisterResponse {
	return dto.RegisterResponse{
		User: serviceUserToDtoUser(user),
	}
}

func serviceUserToDtoUser(user model.User) dto.User {
	return dto.User{
		ID:        user.ID,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	}
}
