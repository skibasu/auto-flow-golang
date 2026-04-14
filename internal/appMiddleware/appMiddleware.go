package appMiddleware

import (
	"reflect"
	"slices"

	"github.com/go-playground/validator/v10"
	"github.com/nyaruka/phonenumbers"
	"github.com/skibasu/auto-flow-api/internal/config"
	"github.com/skibasu/auto-flow-api/internal/utils"
)

type ContextKey string

type UserContext struct {
	Id    string
	Roles []string
}

var UserContextKey = ContextKey("user")

type AppMiddleware struct {
	Validator *validator.Validate
	Config    config.Config
}

func NewAppMiddleware(cfg config.Config) *AppMiddleware {
	v := validator.New(validator.WithRequiredStructEnabled())
	registerValidation(v)

	return &AppMiddleware{Config: cfg, Validator: v}

}

func registerValidation(v *validator.Validate) {

	v.RegisterValidation("password", func(fl validator.FieldLevel) bool {
		password := fl.Field().String()
		return len(password) > 7 && utils.UpperRegex.MatchString(password) && utils.DigitRegex.MatchString(password) && utils.SpecialRegex.MatchString(password)

	})
	v.RegisterValidation("phoneNumber", func(fl validator.FieldLevel) bool {
		num := fl.Field().String()
		if len(num) == 0 || num[0] != '+' {
			return false
		}
		parsed, err := phonenumbers.Parse(num, "")
		if err != nil {
			return false
		}

		return phonenumbers.IsValidNumber(parsed)
	})

	v.RegisterValidation("roles", func(fl validator.FieldLevel) bool {
		field := fl.Field()

		if field.Kind() != reflect.Slice {
			return false
		}

		allowed := config.ROLES

		for i := 0; i < field.Len(); i++ {
			role := field.Index(i).String()

			if !slices.Contains(allowed[:], role) {
				return false
			}
		}

		return true
	})
}
