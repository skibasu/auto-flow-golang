package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/skibasu/auto-flow-api/internal/appMiddleware"
	"github.com/skibasu/auto-flow-api/internal/dto"
	appErrors "github.com/skibasu/auto-flow-api/internal/helpers"
	"github.com/skibasu/auto-flow-api/internal/services"
)

func GetMe(userService *services.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctxUser, ok := r.Context().Value(appMiddleware.UserCtxKey).(appMiddleware.UserContext)
		if !ok {
			appErrors.NewUnauthorized(w, errors.New("invalid user context"))
			return
		}

		if ctxUser.Id == "" {
			appErrors.NewUnauthorized(w, errors.New("missing user id"))
			return
		}

		// 👇 user z bazy
		user, err := userService.GetMe(ctxUser.Id)
		if err != nil {
			appErrors.NewNotFound(w, err)
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
			appErrors.NewInternal(w, err)
			return
		}

		json.NewEncoder(w).Encode(users)
	}
}
