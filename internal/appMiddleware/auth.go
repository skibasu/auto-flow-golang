package appMiddleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/skibasu/auto-flow-api/internal/config"
	"github.com/skibasu/auto-flow-api/internal/helpers"
	"github.com/skibasu/auto-flow-api/internal/jwt"
)

type contextKey string

const UserCtxKey = contextKey("user")

type UserContext struct {
	Id    string
	Roles []string
}

var cfg = config.Load()

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		header := r.Header.Get("Authorization")
		if header == "" {
			helpers.NewUnauthorized(w, errors.New("missing token"), nil)
			return
		}

		parts := strings.Split(header, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			helpers.NewUnauthorized(w, errors.New("invalid token format"), nil)
			return
		}

		tokenStr := parts[1]

		claims, err := jwt.ParseToken(tokenStr, cfg.JWTSecret)
		if err != nil {
			helpers.NewUnauthorized(w, errors.New("invalid token"), nil)
			return
		}
		roles := claims.Roles

		user := UserContext{
			Id:    claims.Sub,
			Roles: roles,
		}

		ctx := context.WithValue(r.Context(), UserCtxKey, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
