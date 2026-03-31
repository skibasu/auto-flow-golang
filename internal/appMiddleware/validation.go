package appMiddleware

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	helpers "github.com/skibasu/auto-flow-api/internal/appErrors"
)

const BodyKey = contextKey("body")

func ValidateRequest[T any](appMiddleware *AppMiddleware, protectData bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			var body T

			decoder := json.NewDecoder(r.Body)
			decoder.DisallowUnknownFields()

			if err := decoder.Decode(&body); err != nil {

				details := map[string]string{}

				var syntaxError *json.SyntaxError
				var unmarshalTypeError *json.UnmarshalTypeError

				switch {

				case errors.As(err, &syntaxError):
					details["json"] = "invalid JSON syntax"

				case errors.As(err, &unmarshalTypeError):
					details[unmarshalTypeError.Field] = unmarshalTypeError.Field + ": invalid type"

				case strings.HasPrefix(err.Error(), "json: unknown field"):
					field := strings.TrimPrefix(err.Error(), "json: unknown field ")
					details["json"] = "unknown field " + field

				case errors.Is(err, io.EOF):
					details["json"] = "empty body"

				default:
					details["json"] = "invalid request"
				}
				if protectData {
					helpers.NewUnauthorized(w, errors.New("invalid credentials"), nil)
					return
				} else {
					helpers.NewBadRequest(w, errors.New("invalid request"), &details)
					return
				}

			}

			if err := appMiddleware.Validator.Struct(body); err != nil {
				var validateErrs validator.ValidationErrors

				if errors.As(err, &validateErrs) {
					details := map[string]string{}

					for _, e := range validateErrs {
						details[strings.ToLower(e.Field())] = e.Tag()
					}
					if protectData {
						helpers.NewUnauthorized(w, errors.New("invalid credentials"), nil)
						return
					} else {
						helpers.NewBadRequest(w, errors.New("validation error"), &details)
						return
					}

				}

				helpers.NewBadRequest(w, errors.New("validation error"), nil)
				return
			}

			ctx := context.WithValue(r.Context(), BodyKey, body)

			next.ServeHTTP(w, r.WithContext(ctx))

		})
	}
}
func GetValidatedBody[T any](r *http.Request) T {
	val := r.Context().Value(BodyKey)
	if val == nil {
		panic("GetValidatedBody called without ValidateRequest middleware")
	}
	return val.(T)
}
