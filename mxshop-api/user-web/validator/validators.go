package validator

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func ValidateMobile(fl validator.FieldLevel) bool {
	mobile := fl.Field().String()
	ok, _ := regexp.MatchString(`^1[3-9]\d{9}$`, mobile)
	return ok
}
