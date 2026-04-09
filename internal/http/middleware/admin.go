package middleware

import (
	"net/http"

	"github.com/soundmarket/backend/internal/domain"
	"github.com/soundmarket/backend/internal/http/response"
)

func RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := CurrentUser(r)
		if user.Role != domain.RoleAdmin {
			response.Error(w, http.StatusForbidden, "forbidden")
			return
		}
		next.ServeHTTP(w, r)
	})
}
