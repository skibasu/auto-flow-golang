package appMiddleware

import (
	"github.com/go-playground/validator/v10"
	"github.com/nyaruka/phonenumbers"
)

func RegisterValidation() {
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
		return true

	})
}
