package service

import (
	"github.com/go-playground/validator/v10"
)

type (
	IError struct {
		Error       bool        `json:"error"`
		FailedField string      `json:"failed_field"`
		Tag         string      `json:"tag"`
		Value       interface{} `json:"value"`
	}

	XValidator struct {
		Validator *validator.Validate
	}

	GlobalErrorHandlerResp struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}
)

var validating = validator.New()

func (v XValidator) Validate(data interface{}) []IError {
	validationErrors := []IError{}

	errs := validating.Struct(data)
	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			var elem IError
			elem.FailedField = err.Field()
			elem.Tag = err.Tag()
			elem.Value = err.Value()
			elem.Error = true

			validationErrors = append(validationErrors, elem)
		}
	}

	return validationErrors
}
