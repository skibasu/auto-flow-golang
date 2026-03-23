package appErrors

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Message string             `json:"message"`
	Status  int                `json:"status"`
	Code    string             `json:"code"`
	Details *map[string]string `json:"details,omitempty"`
}

func errorResponse(w http.ResponseWriter, err ErrorResponse) {

	w.WriteHeader(err.Status)
	json.NewEncoder(w).Encode(err)
}
func NewNotFound(w http.ResponseWriter, err error, details *map[string]string) {
	res := ErrorResponse{
		Message: err.Error(),
		Code:    "NOT_FOUND",
		Status:  http.StatusNotFound,
		Details: details,
	}
	errorResponse(w, res)
}

func NewBadRequest(w http.ResponseWriter, err error, details *map[string]string) {

	res := ErrorResponse{
		Message: err.Error(),
		Code:    "BAD_REQUEST",
		Status:  http.StatusBadRequest,
		Details: details,
	}
	errorResponse(w, res)
}

func NewUnauthorized(w http.ResponseWriter, err error, details *map[string]string) {
	res := ErrorResponse{
		Message: err.Error(),
		Code:    "UNAUTHORIZED",
		Status:  http.StatusUnauthorized,
		Details: details,
	}
	errorResponse(w, res)
}

func NewForbidden(w http.ResponseWriter, err error, details *map[string]string) {
	res := ErrorResponse{
		Message: err.Error(),
		Code:    "FORBIDDEN",
		Status:  http.StatusForbidden,
		Details: details,
	}
	errorResponse(w, res)
}

func NewInternal(w http.ResponseWriter, err error, details *map[string]string) {
	res := ErrorResponse{
		Message: err.Error(),
		Code:    "INTERNAL_ERROR",
		Status:  http.StatusInternalServerError,
		Details: details,
	}
	errorResponse(w, res)
}
