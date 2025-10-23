package response

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation"
)

var validationErrs validation.Errors

type fieldError struct {
	Field string `json:"field"`
	Msg   string `json:"description"`
}

func extractFieldError(result *[]fieldError, parentField string, v any) {
	for f, err := range v.(validation.Errors) {
		field := fmt.Sprintf("%s[%s]", parentField, f)
		switch err := err.(type) {
		case validation.Errors:
			extractFieldError(result, field, err)
		default:
			*result = append(*result, fieldError{
				Field: field,
				Msg:   err.Error(),
			})
		}
	}
}

// convertValidationErrors transform validation errors to slice map.
func convertValidationErrors(v validation.Errors) []fieldError {
	var errs []fieldError
	for field, err := range v {
		switch err := err.(type) {
		case validation.Errors:
			extractFieldError(&errs, field, err)
		default:
			errs = append(errs, fieldError{
				Field: field,
				Msg:   err.Error(),
			})
		}
	}

	return errs
}
