package appMiddleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/skibasu/auto-flow-api/internal/appErrors"
	"github.com/skibasu/auto-flow-api/internal/jwt"
)

func (m *AppMiddleware) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		header := r.Header.Get("Authorization")
		if header == "" {
			appErrors.NewUnauthorized(w, errors.New("missing token"), nil)
			return
		}

		parts := strings.Split(header, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			appErrors.NewUnauthorized(w, errors.New("invalid token format"), nil)
			return
		}

		tokenStr := parts[1]

		claims, err := jwt.ParseToken(tokenStr, m.Config.Secret)
		if err != nil {
			appErrors.NewUnauthorized(w, errors.New("invalid token"), nil)
			return
		}
		roles := claims.Roles

		user := UserContext{
			Id:    claims.Sub,
			Roles: roles,
		}

		ctx := context.WithValue(r.Context(), m.UserCtxKey, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
