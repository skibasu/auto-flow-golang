package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/skibasu/auto-flow-api/internal/appMiddleware"
	"github.com/skibasu/auto-flow-api/internal/dto"
	appErrors "github.com/skibasu/auto-flow-api/internal/helpers"
	"github.com/skibasu/auto-flow-api/internal/services"
)

func GetMe(userService *services.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctxUser, ok := r.Context().Value(appMiddleware.UserCtxKey).(appMiddleware.UserContext)
		if !ok {
			appErrors.NewUnauthorized(w, errors.New("invalid user context"), nil)
			return
		}

		if ctxUser.Id == "" {
			appErrors.NewUnauthorized(w, errors.New("missing user id"), nil)
			return
		}

		// 👇 user z bazy
		user, err := userService.GetMe(ctxUser.Id)
		if err != nil {
			appErrors.NewNotFound(w, err, nil)
			return
		}

		json.NewEncoder(w).Encode(user)
	}
}

func GetUsers(userService *services.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		urlQuery := r.URL.Query()

		filter := dto.UsersFilterRequest{
			Email:       urlQuery.Get("email"),
			FirstName:   urlQuery.Get("firstName"),
			LastName:    urlQuery.Get("lastName"),
			PhoneNumber: urlQuery.Get("phoneNumber"),
		}
		if roles := urlQuery.Get("roles"); roles != "" {
			filter.Roles = strings.Split(roles, ",")
		}

		users, err := userService.GetUsers(filter)
		if err != nil {
			appErrors.NewInternal(w, err, nil)
			return
		}

		json.NewEncoder(w).Encode(users)
	}
}

func CreateUser(userService *services.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := appMiddleware.GetValidatedBody[dto.UserRequest](r)

		user, err := userService.CreateUser(req)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") {
				appErrors.NewConflict(w, err, nil)
				return
			}
			appErrors.NewInternal(w, err, nil)
			return
		}
		json.NewEncoder(w).Encode(user)
	}
}

func DeleteUser(userService *services.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		id := chi.URLParam(r, "id")

		if id == "" {
			appErrors.NewBadRequest(w, errors.New("missing id"), nil)
			return
		}

		err := userService.DeleteUser(id)

		if err != nil {
			if strings.Contains(err.Error(), "user not found") {
				appErrors.NewNotFound(w, err, nil)
				return
			}
			appErrors.NewInternal(w, err, nil)
			return
		}

		w.WriteHeader(http.StatusNoContent) // 204

	}
}

func UpdateUser(userService *services.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		id := chi.URLParam(r, "id")

		if id == "" {
			appErrors.NewBadRequest(w, errors.New("missing id"), nil)
			return
		}

		req := appMiddleware.GetValidatedBody[dto.UpdateUserRequest](r)
		user, err := userService.UpdateUser(id, req)
		if err != nil {

			appErrors.NewInternal(w, err, nil)
			return
		}
		json.NewEncoder(w).Encode(user)

	}
}
