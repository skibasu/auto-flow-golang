package appMiddleware

import (
	"errors"
	"net/http"
	"slices"

	helpers "github.com/skibasu/auto-flow-api/internal/appErrors"
)

func hasRequiredRole(userRoles []string, requiredRoles []string) bool {
	for _, u := range userRoles {
		if slices.Contains(requiredRoles, u) {
			return true
		}
	}
	return false
}

func (m *AppMiddleware) RequireRole(roles []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			user, ok := r.Context().Value(UserContextKey).(UserContext)
			if !ok {
				helpers.NewUnauthorized(w, errors.New("unauthorized"), nil)

				return
			}

			if !hasRequiredRole(user.Roles, roles) {
				helpers.NewForbidden(w, errors.New("forbidden"), nil)

				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
