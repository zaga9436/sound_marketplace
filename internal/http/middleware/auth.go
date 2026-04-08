package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/soundmarket/backend/internal/auth"
	"github.com/soundmarket/backend/internal/domain"
	"github.com/soundmarket/backend/internal/http/response"
)

type contextKey string

const userContextKey contextKey = "current_user"

func RequireAuth(jwtManager *auth.JWTManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if !strings.HasPrefix(header, "Bearer ") {
				response.Error(w, http.StatusUnauthorized, "missing bearer token")
				return
			}
			claims, err := jwtManager.Parse(strings.TrimPrefix(header, "Bearer "))
			if err != nil {
				response.Error(w, http.StatusUnauthorized, "invalid token")
				return
			}
			ctx := context.WithValue(r.Context(), userContextKey, domain.User{
				ID:   claims.UserID,
				Role: domain.Role(claims.Role),
			})
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func CurrentUser(r *http.Request) domain.User {
	user, _ := r.Context().Value(userContextKey).(domain.User)
	return user
}
