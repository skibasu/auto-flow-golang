package appErrors

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
	Code    string `json:"code"`
}

func errorResponse(w http.ResponseWriter, status int, err ErrorResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(err)
}
func NewNotFound(w http.ResponseWriter, err error) {
	res := ErrorResponse{
		Message: err.Error(),
		Code:    "NOT_FOUND",
		Status:  404,
	}
	errorResponse(w, 404, res)
}

func NewBadRequest(w http.ResponseWriter, err error) {

	res := ErrorResponse{
		Message: err.Error(),
		Code:    "BAD_REQUEST",
		Status:  400,
	}
	errorResponse(w, 400, res)
}

func NewUnauthorized(w http.ResponseWriter, err error) {
	res := ErrorResponse{
		Message: err.Error(),
		Code:    "UNAUTHORIZED",
		Status:  401,
	}
	errorResponse(w, 401, res)
}

func NewForbidden(w http.ResponseWriter, err error) {
	res := ErrorResponse{
		Message: err.Error(),
		Code:    "FORBIDDEN",
		Status:  403,
	}
	errorResponse(w, 403, res)
}

func NewInternal(w http.ResponseWriter, err error) {
	res := ErrorResponse{
		Message: err.Error(),
		Code:    "INTERNAL_ERROR",
		Status:  500,
	}
	errorResponse(w, 500, res)
}
