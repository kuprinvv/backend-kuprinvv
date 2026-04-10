package middleware

import (
	"net/http"
	"net/http/httptest"
	"test-backend-1-kuprinvv/internal/model"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testSecret = "test-secret-key"

type staticJWTConfig struct{}

func (staticJWTConfig) Token() string { return testSecret }

// makeToken создаёт подписанный JWT с заданными claims.
func makeToken(userID uuid.UUID, role string, expired bool) string {
	exp := time.Now().Add(24 * time.Hour)
	if expired {
		exp = time.Now().Add(-1 * time.Hour)
	}
	claims := model.AuthClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, _ := token.SignedString([]byte(testSecret))
	return signed
}

// nextHandler — заглушка, которая пишет 200 и сохраняет UserContext.
func nextHandler(captured *UserContext) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u := GetUser(r.Context())
		if captured != nil {
			*captured = u
		}
		w.WriteHeader(http.StatusOK)
	})
}

func TestAuthMiddleware(t *testing.T) {
	userID := uuid.New()
	cfg := staticJWTConfig{}

	tests := []struct {
		name       string
		header     string
		wantStatus int
		wantRole   string
	}{
		{
			name:       "валидный токен admin",
			header:     "Bearer " + makeToken(userID, "admin", false),
			wantStatus: http.StatusOK,
			wantRole:   "admin",
		},
		{
			name:       "валидный токен user",
			header:     "Bearer " + makeToken(userID, "user", false),
			wantStatus: http.StatusOK,
			wantRole:   "user",
		},
		{
			name:       "заголовок отсутствует",
			header:     "",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "нет префикса Bearer",
			header:     makeToken(userID, "user", false),
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "невалидная подпись",
			header:     "Bearer eyJhbGciOiJIUzI1NiJ9.eyJ1c2VyX2lkIjoiMTIzIn0.wrongsig",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "истёкший токен",
			header:     "Bearer " + makeToken(userID, "user", true),
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var captured UserContext
			handler := AuthMiddleware(cfg)(nextHandler(&captured))

			r := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.header != "" {
				r.Header.Set("Authorization", tt.header)
			}
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, r)
			assert.Equal(t, tt.wantStatus, w.Code)

			if tt.wantStatus == http.StatusOK {
				assert.Equal(t, tt.wantRole, captured.Role)
				assert.Equal(t, userID, captured.UserID)
			}
		})
	}
}

func TestRequireRole(t *testing.T) {
	userID := uuid.New()

	tests := []struct {
		name         string
		requiredRole string
		userRole     string
		wantStatus   int
	}{
		{
			name:         "роль совпадает — admin",
			requiredRole: "admin",
			userRole:     "admin",
			wantStatus:   http.StatusOK,
		},
		{
			name:         "роль совпадает — user",
			requiredRole: "user",
			userRole:     "user",
			wantStatus:   http.StatusOK,
		},
		{
			name:         "роль не совпадает",
			requiredRole: "admin",
			userRole:     "user",
			wantStatus:   http.StatusForbidden,
		},
		{
			name:         "пользователь пытается получить ресурс администратора",
			requiredRole: "admin",
			userRole:     "user",
			wantStatus:   http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := RequireRole(tt.requiredRole)(nextHandler(nil))

			r := httptest.NewRequest(http.MethodGet, "/", nil)
			ctx := WithUser(r.Context(), UserContext{UserID: userID, Role: tt.userRole})
			r = r.WithContext(ctx)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, r)
			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestWithUser(t *testing.T) {
	userID := uuid.New()
	user := UserContext{UserID: userID, Role: "admin"}

	ctx := WithUser(t.Context(), user)
	got := GetUser(ctx)

	require.Equal(t, user.UserID, got.UserID)
	require.Equal(t, user.Role, got.Role)
}
