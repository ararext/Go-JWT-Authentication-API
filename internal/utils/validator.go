package utils

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func ValidateStruct(s interface{}) error {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	var messages []string
	for _, fieldErr := range err.(validator.ValidationErrors) {
		messages = append(messages, formatFieldError(fieldErr))
	}

	return errors.New(strings.Join(messages, "; "))
}

func formatFieldError(fe validator.FieldError) string {
	field := fe.Field()

	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "email":
		return fmt.Sprintf("%s must be a valid email address", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", field, fe.Param())
	default:
		return fmt.Sprintf("%s is invalid", field)
	}
}
