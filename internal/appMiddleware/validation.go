package appMiddleware

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	appErrors "github.com/skibasu/auto-flow-api/internal/helpers"
)

const BodyKey = contextKey("body")

var validate = validator.New(validator.WithRequiredStructEnabled())

func ValidateRequest[T any](hideDetails bool) func(http.Handler) http.Handler {
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
				if hideDetails {
					appErrors.NewUnauthorized(w, errors.New("invalid credentials"), nil)
					return
				} else {
					appErrors.NewBadRequest(w, errors.New("invalid request"), &details)
					return
				}

			}

			if err := validate.Struct(body); err != nil {
				var validateErrs validator.ValidationErrors

				if errors.As(err, &validateErrs) {
					details := map[string]string{}

					for _, e := range validateErrs {
						details[strings.ToLower(e.Field())] = e.Tag()
					}
					if hideDetails {
						appErrors.NewUnauthorized(w, errors.New("invalid credentials"), nil)
						return
					} else {
						appErrors.NewBadRequest(w, errors.New("validation error"), &details)
						return
					}

				}

				appErrors.NewBadRequest(w, err, nil)
				return
			}

			ctx := context.WithValue(r.Context(), BodyKey, body)

			next.ServeHTTP(w, r.WithContext(ctx))

		})
	}
}

func GetValidatedBody[T any](r *http.Request) T {
	return r.Context().Value(BodyKey).(T)
}
