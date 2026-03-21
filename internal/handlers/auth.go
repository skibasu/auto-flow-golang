package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/skibasu/auto-flow-api/internal/dto"
	appErrors "github.com/skibasu/auto-flow-api/internal/helpers"
	"github.com/skibasu/auto-flow-api/internal/services"
)

func Auth(authService *services.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var req dto.LoginRequest

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			appErrors.NewBadRequest(w, err)
			return
		}

		access, refresh, err := authService.Login(req.Email, req.Password)
		if err != nil {
			appErrors.NewUnauthorized(w, err)
			return
		}

		// 🍪 refresh token
		http.SetCookie(w, &http.Cookie{
			Name:     "refreshToken",
			Value:    refresh,
			HttpOnly: true,
			Path:     "/",
		})

		// 📦 response
		json.NewEncoder(w).Encode(map[string]string{
			"accessToken": access,
		})
	}
}

func RefreshToken(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refreshToken")
	if err != nil {
		appErrors.NewUnauthorized(w, errors.New("missing refresh token"))
		return
	}

	service := services.AuthService{}

	access, refresh, err := service.Refresh(cookie.Value)
	if err != nil {
		appErrors.NewUnauthorized(w, errors.New("invalid refresh token"))

		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refreshToken",
		Value:    refresh,
		HttpOnly: true,
		Path:     "/",
	})

	json.NewEncoder(w).Encode(map[string]string{
		"accessToken": access,
	})
}
