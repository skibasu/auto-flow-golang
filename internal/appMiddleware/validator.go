package appMiddleware

import (
	"reflect"
	"regexp"
	"slices"

	"github.com/go-playground/validator/v10"
	"github.com/nyaruka/phonenumbers"
	"github.com/skibasu/auto-flow-api/internal/config"
)

var (
	upperRegex   = regexp.MustCompile(`[A-Z]`)
	digitRegex   = regexp.MustCompile(`\d`)
	specialRegex = regexp.MustCompile(`[!@#$%^&*()_\-+=<>?{}\[\]~]`)
)

func RegisterValidation() {

	validate.RegisterValidation("password", func(fl validator.FieldLevel) bool {
		password := fl.Field().String()
		return len(password) > 7 && upperRegex.MatchString(password) && digitRegex.MatchString(password) && specialRegex.MatchString(password)

	})
	validate.RegisterValidation("phoneNumber", func(fl validator.FieldLevel) bool {
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

	validate.RegisterValidation("roles", func(fl validator.FieldLevel) bool {
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
