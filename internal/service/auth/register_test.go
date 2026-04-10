package auth

import (
	"context"
	"errors"
	"test-backend-1-kuprinvv/internal/mocks"
	"test-backend-1-kuprinvv/internal/model"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type testJWTConfig struct{}

func (c *testJWTConfig) Token() string { return "test-secret-key" }

func TestRegister(t *testing.T) {
	ctx := context.Background()
	dbErr := errors.New("db error")

	tests := []struct {
		name     string
		email    string
		password string
		role     string
		setup    func(repo *mocks.MockUserRepository)
		wantErr  error
		wantUser bool
	}{
		{
			name:     "успешная регистрация",
			email:    "user@example.com",
			password: "pass",
			role:     "user",
			setup:    func(repo *mocks.MockUserRepository) { repo.EXPECT().Create(ctx, gomock.Any()).Return(nil) },
			wantUser: true,
		},
		{
			name:     "email уже занят",
			email:    "taken@example.com",
			password: "pass",
			role:     "user",
			setup: func(repo *mocks.MockUserRepository) {
				repo.EXPECT().Create(ctx, gomock.Any()).Return(&pgconn.PgError{Code: "23505"})
			},
			wantErr: model.ErrUserAlreadyExists,
		},
		{
			name:     "ошибка репозитория",
			email:    "user@example.com",
			password: "pass",
			role:     "user",
			setup:    func(repo *mocks.MockUserRepository) { repo.EXPECT().Create(ctx, gomock.Any()).Return(dbErr) },
			wantErr:  dbErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mocks.NewMockUserRepository(ctrl)
			tt.setup(repo)

			svc := NewAuthService(&testJWTConfig{}, repo)
			user, err := svc.Register(ctx, tt.email, tt.password, tt.role)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				assert.Empty(t, user.ID)
			} else {
				require.NoError(t, err)
				assert.NotEmpty(t, user.ID)
				assert.Equal(t, tt.email, user.Email)
			}
		})
	}
}

func TestDummyLogin(t *testing.T) {
	ctx := context.Background()
	adminID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	dbErr := errors.New("db error")

	tests := []struct {
		name      string
		userID    uuid.UUID
		role      string
		setup     func(repo *mocks.MockUserRepository)
		wantErr   bool
		wantToken bool
	}{
		{
			name:      "успешно: первый вызов создаёт пользователя",
			userID:    adminID,
			role:      "admin",
			setup:     func(repo *mocks.MockUserRepository) { repo.EXPECT().Create(ctx, gomock.Any()).Return(nil) },
			wantToken: true,
		},
		{
			name:   "идемпотентность: нарушение уникальности игнорируется, токен всё равно возвращается",
			userID: adminID,
			role:   "admin",
			setup: func(repo *mocks.MockUserRepository) {
				repo.EXPECT().Create(ctx, gomock.Any()).Return(&pgconn.PgError{Code: "23505"})
			},
			wantToken: true,
		},
		{
			name:    "ошибка репозитория передаётся наверх",
			userID:  adminID,
			role:    "admin",
			setup:   func(repo *mocks.MockUserRepository) { repo.EXPECT().Create(ctx, gomock.Any()).Return(dbErr) },
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mocks.NewMockUserRepository(ctrl)
			tt.setup(repo)

			svc := NewAuthService(&testJWTConfig{}, repo)
			token, err := svc.DummyLogin(ctx, tt.userID, tt.role)

			if tt.wantErr {
				require.Error(t, err)
				assert.Empty(t, token)
			} else {
				require.NoError(t, err)
				assert.NotEmpty(t, token)
			}
		})
	}
}
