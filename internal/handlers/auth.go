package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/skibasu/auto-flow-api/internal/appMiddleware"
	"github.com/skibasu/auto-flow-api/internal/dto"
	appErrors "github.com/skibasu/auto-flow-api/internal/helpers"
	"github.com/skibasu/auto-flow-api/internal/services"
)

func (h *Handler) Auth(authService *services.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := appMiddleware.GetValidatedBody[dto.Credentials](r)

		access, refresh, err := authService.Login(req.Email, req.Password)
		if err != nil {
			appErrors.NewUnauthorized(w, err, nil)
			return
		}

		// 🍪 refresh token
		http.SetCookie(w, &http.Cookie{
			Name:     "refreshToken",
			Value:    refresh,
			HttpOnly: true,
			Path:     "/",
			Secure:   true,
			SameSite: http.SameSiteNoneMode,
		})

		// 📦 response
		json.NewEncoder(w).Encode(map[string]string{
			"accessToken": access,
		})
	}
}

func (h *Handler) RefreshToken(authService *services.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		{
			cookie, err := r.Cookie("refreshToken")
			if err != nil {
				appErrors.NewUnauthorized(w, errors.New("missing refresh token"), nil)
				return
			}

			access, refresh, err := authService.Refresh(cookie.Value)
			if err != nil {
				appErrors.NewUnauthorized(w, errors.New("invalid refresh token"), nil)

				return
			}

			http.SetCookie(w, &http.Cookie{
				Name:     "refreshToken",
				Value:    refresh,
				HttpOnly: true,
				Path:     "/",
				Secure:   true,
				SameSite: http.SameSiteNoneMode,
			})

			json.NewEncoder(w).Encode(map[string]string{
				"accessToken": access,
			})
		}
	}
}
