package auth

import (
	"context"
	"errors"
	"test-backend-1-kuprinvv/internal/mocks"
	"test-backend-1-kuprinvv/internal/model"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

func TestLogin(t *testing.T) {
	ctx := context.Background()

	hash, err := bcrypt.GenerateFromPassword([]byte("correct-password"), bcrypt.MinCost)
	require.NoError(t, err)

	userID := uuid.New()
	storedUser := &model.User{
		ID:       userID,
		Email:    "user@example.com",
		Password: string(hash),
		Role:     "user",
	}

	tests := []struct {
		name      string
		email     string
		password  string
		setup     func(repo *mocks.MockUserRepository)
		wantErr   error
		wantToken bool
	}{
		{
			name:     "успешный вход",
			email:    "user@example.com",
			password: "correct-password",
			setup: func(repo *mocks.MockUserRepository) {
				repo.EXPECT().GetByEmail(ctx, "user@example.com").Return(storedUser, nil)
			},
			wantToken: true,
		},
		{
			name:     "пользователь не найден",
			email:    "unknown@example.com",
			password: "any",
			setup: func(repo *mocks.MockUserRepository) {
				repo.EXPECT().GetByEmail(ctx, "unknown@example.com").Return(nil, pgx.ErrNoRows)
			},
			wantErr: model.ErrInvalidCredentials,
		},
		{
			name:     "неверный пароль",
			email:    "user@example.com",
			password: "wrong-password",
			setup: func(repo *mocks.MockUserRepository) {
				repo.EXPECT().GetByEmail(ctx, "user@example.com").Return(storedUser, nil)
			},
			wantErr: model.ErrInvalidCredentials,
		},
		{
			name:     "ошибка репозитория",
			email:    "user@example.com",
			password: "any",
			setup: func(repo *mocks.MockUserRepository) {
				repo.EXPECT().GetByEmail(ctx, "user@example.com").Return(nil, errors.New("db error"))
			},
			wantErr: errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mocks.NewMockUserRepository(ctrl)
			tt.setup(repo)

			svc := NewAuthService(&testJWTConfig{}, repo)
			token, err := svc.Login(ctx, tt.email, tt.password)

			if tt.wantErr != nil {
				require.Error(t, err)
				assert.Empty(t, token)
			} else {
				require.NoError(t, err)
				assert.NotEmpty(t, token)
			}
		})
	}
}
