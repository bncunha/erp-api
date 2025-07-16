package validator

import "github.com/go-playground/validator/v10"

func Validate(data any) error {
	return validator.New().Struct(data)
}