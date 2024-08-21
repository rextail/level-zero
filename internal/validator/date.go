package validator

import (
	"github.com/go-playground/validator/v10"
	"time"
)

func ValidateDate(fl validator.FieldLevel) bool {
	t := fl.Field().String()

	_, err := time.Parse(time.RFC3339, t)

	return err == nil
}
