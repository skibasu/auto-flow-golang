package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/skibasu/auto-flow-api/internal/appMiddleware"
)

func GetMeAdmin(w http.ResponseWriter, r *http.Request) {

	user, ok := r.Context().Value(appMiddleware.UserCtxKey).(appMiddleware.UserContext)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(user)
}
