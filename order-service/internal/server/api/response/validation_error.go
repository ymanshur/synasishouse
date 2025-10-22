package response

import (
	"github.com/go-playground/validator/v10"
)

type validationError struct {
	Field string `json:"field"`
	Msg   string `json:"description"`
}

// validationMsgFromFieldError translates error message from field
func validationMsgFromFieldError(field validator.FieldError) string {
	switch field.Tag() {
	case "required":
		return "is required"
	case "uuid":
		return "must be uuid"
	case "datetime":
		return "must be 2006-01-02T15:04:05Z07:00 format"
	}
	return field.Error()
}
