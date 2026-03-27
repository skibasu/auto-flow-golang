package appMiddleware

import (
	"errors"
	"net/http"
	"slices"

	helpers "github.com/skibasu/auto-flow-api/internal/appErrors"
)

func HasRequiredRole(userRoles []string, requiredRoles []string) bool {
	for _, u := range userRoles {
		if slices.Contains(requiredRoles, u) {
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
				helpers.NewUnauthorized(w, errors.New("unauthorized"), nil)

				return
			}

			if !HasRequiredRole(user.Roles, roles) {
				helpers.NewForbidden(w, errors.New("forbidden"), nil)

				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
