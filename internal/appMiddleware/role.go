package appMiddleware

import (
	"errors"
	"net/http"
	"slices"

	appErrors "github.com/skibasu/auto-flow-api/internal/helpers"
)

func hasRequiredRole(userRoles []string, requiredRoles []string) bool {
	for _, ur := range userRoles {
		if slices.Contains(requiredRoles, ur) {
			return true
		}
	}
	return false
}

func RequireRole(roles []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			user, ok := r.Context().Value(UserCtxKey).(UserContext)
			if !ok {
				appErrors.NewUnauthorized(w, errors.New("unauthorized"))

				return
			}

			if !hasRequiredRole(user.Roles, roles) {
				appErrors.NewForbidden(w, errors.New("forbidden"))

				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
