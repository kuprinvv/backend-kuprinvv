package middleware

import (
	"context"
	"net/http"
	"strings"
	"test-backend-1-kuprinvv/internal/config"
	"test-backend-1-kuprinvv/internal/handler"
	"test-backend-1-kuprinvv/internal/model"
	"test-backend-1-kuprinvv/pkg/httpx"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type contextKey string

const userContextKey = contextKey("user")

type UserContext struct {
	UserID uuid.UUID
	Role   string
}

func AuthMiddleware(cfg config.JWTConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				httpx.Error(w, handler.ErrUnauthorized, "invalid token", http.StatusUnauthorized)
				return
			}

			if !strings.HasPrefix(authHeader, "Bearer ") {
				httpx.Error(w, handler.ErrUnauthorized, "invalid token", http.StatusUnauthorized)
				return
			}
			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

			token, err := jwt.ParseWithClaims(tokenStr, &model.AuthClaims{}, func(token *jwt.Token) (interface{}, error) {
				return []byte(cfg.Token()), nil
			})
			if err != nil {
				httpx.Error(w, handler.ErrUnauthorized, "invalid token", http.StatusUnauthorized)
				return
			}
			if !token.Valid {
				httpx.Error(w, handler.ErrUnauthorized, "token expired", http.StatusUnauthorized)
				return
			}

			claims, ok := token.Claims.(*model.AuthClaims)
			if !ok {
				httpx.Error(w, handler.ErrUnauthorized, "invalid token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), userContextKey, UserContext{
				UserID: claims.UserID,
				Role:   claims.Role,
			})

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUser(ctx context.Context) UserContext {
	return ctx.Value(userContextKey).(UserContext)
}

// WithUser возвращает копию ctx с установленным UserContext.
// Предназначен для использования в тестах.
func WithUser(ctx context.Context, user UserContext) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}
