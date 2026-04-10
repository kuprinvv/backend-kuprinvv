package middleware

import (
	"net/http"
	"test-backend-1-kuprinvv/internal/handler"
	"test-backend-1-kuprinvv/pkg/httpx"
)

func RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := GetUser(r.Context())

			if user.Role != role {
				httpx.Error(w, handler.ErrForbidden, "invalid role", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
