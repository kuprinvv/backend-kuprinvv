package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"test-backend-1-kuprinvv/internal/mocks"
	"test-backend-1-kuprinvv/internal/model"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestDummyLogin(t *testing.T) {
	tests := []struct {
		name       string
		body       any
		setup      func(svc *mocks.MockAuthService)
		wantStatus int
	}{
		{
			name: "успешный вход admin",
			body: map[string]any{"role": "admin"},
			setup: func(svc *mocks.MockAuthService) {
				svc.EXPECT().
					DummyLogin(gomock.Any(), gomock.Any(), "admin").
					Return("test-token", nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "успешный вход user",
			body: map[string]any{"role": "user"},
			setup: func(svc *mocks.MockAuthService) {
				svc.EXPECT().
					DummyLogin(gomock.Any(), gomock.Any(), "user").
					Return("test-token", nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "невалидный JSON",
			body:       "not json",
			setup:      func(svc *mocks.MockAuthService) {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "пустая роль",
			body:       map[string]any{"role": ""},
			setup:      func(svc *mocks.MockAuthService) {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "некорректная роль",
			body:       map[string]any{"role": "superuser"},
			setup:      func(svc *mocks.MockAuthService) {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "внутренняя ошибка",
			body: map[string]any{"role": "user"},
			setup: func(svc *mocks.MockAuthService) {
				svc.EXPECT().
					DummyLogin(gomock.Any(), gomock.Any(), "user").
					Return("", errors.New("db error"))
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			svc := mocks.NewMockAuthService(ctrl)
			tt.setup(svc)

			h := NewAuthHandler(svc)

			rawBody, _ := json.Marshal(tt.body)
			r := httptest.NewRequest(http.MethodPost, "/dummyLogin", bytes.NewReader(rawBody))
			r.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			h.DummyLogin(w, r)
			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestLogin(t *testing.T) {
	tests := []struct {
		name       string
		body       any
		setup      func(svc *mocks.MockAuthService)
		wantStatus int
	}{
		{
			name: "успешный вход",
			body: map[string]any{"email": "user@example.com", "password": "secret123"},
			setup: func(svc *mocks.MockAuthService) {
				svc.EXPECT().
					Login(gomock.Any(), "user@example.com", "secret123").
					Return("test-token", nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "невалидный JSON",
			body:       "not json",
			setup:      func(svc *mocks.MockAuthService) {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "пустой email",
			body:       map[string]any{"email": "", "password": "secret123"},
			setup:      func(svc *mocks.MockAuthService) {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "пароль < 6 символов",
			body:       map[string]any{"email": "user@example.com", "password": "abc"},
			setup:      func(svc *mocks.MockAuthService) {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "неверный пароль",
			body: map[string]any{"email": "user@example.com", "password": "wrongpass"},
			setup: func(svc *mocks.MockAuthService) {
				svc.EXPECT().
					Login(gomock.Any(), "user@example.com", "wrongpass").
					Return("", model.ErrInvalidCredentials)
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "внутренняя ошибка",
			body: map[string]any{"email": "user@example.com", "password": "secret123"},
			setup: func(svc *mocks.MockAuthService) {
				svc.EXPECT().
					Login(gomock.Any(), "user@example.com", "secret123").
					Return("", errors.New("db error"))
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			svc := mocks.NewMockAuthService(ctrl)
			tt.setup(svc)

			h := NewAuthHandler(svc)

			rawBody, _ := json.Marshal(tt.body)
			r := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(rawBody))
			r.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			h.Login(w, r)
			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestRegister(t *testing.T) {
	tests := []struct {
		name       string
		body       any
		setup      func(svc *mocks.MockAuthService)
		wantStatus int
	}{
		{
			name: "успешная регистрация",
			body: map[string]any{"email": "new@example.com", "password": "secret123", "role": "user"},
			setup: func(svc *mocks.MockAuthService) {
				svc.EXPECT().
					Register(gomock.Any(), "new@example.com", "secret123", "user").
					Return(model.User{Email: "new@example.com", Role: "user"}, nil)
			},
			wantStatus: http.StatusCreated,
		},
		{
			name:       "невалидный JSON",
			body:       "not json",
			setup:      func(svc *mocks.MockAuthService) {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "некорректный email",
			body:       map[string]any{"email": "not-email", "password": "secret123", "role": "user"},
			setup:      func(svc *mocks.MockAuthService) {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "пользователь уже существует",
			body: map[string]any{"email": "existing@example.com", "password": "secret123", "role": "user"},
			setup: func(svc *mocks.MockAuthService) {
				svc.EXPECT().
					Register(gomock.Any(), "existing@example.com", "secret123", "user").
					Return(model.User{}, model.ErrUserAlreadyExists)
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "внутренняя ошибка",
			body: map[string]any{"email": "new@example.com", "password": "secret123", "role": "admin"},
			setup: func(svc *mocks.MockAuthService) {
				svc.EXPECT().
					Register(gomock.Any(), "new@example.com", "secret123", "admin").
					Return(model.User{}, errors.New("db error"))
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			svc := mocks.NewMockAuthService(ctrl)
			tt.setup(svc)

			h := NewAuthHandler(svc)

			rawBody, _ := json.Marshal(tt.body)
			r := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(rawBody))
			r.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			h.Register(w, r)
			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}
