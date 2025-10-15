package utils

import (
	"api/internal/model"
	"fmt"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func ValidateStruct(s interface{}) ([]model.ErrorDetails, string, error) {
	err := validate.Struct(s)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		var details []model.ErrorDetails
		for _, e := range errors {
			details = append(details, model.ErrorDetails{
				Field:   e.Field(),
				Message: ConvertTagToMessage(e.Tag(), e.Param()),
			})
		}
		return details, "Validation error", err
	}
	return nil, "", nil
}

func ConvertTagToMessage(tag string, param string) string {
	switch tag {
	case "required":
		return "This field is required."
	case "email":
		return "Invalid email format."
	case "min":
		return fmt.Sprintf("This field must have at least %s characters.", param)
	case "max":
		return fmt.Sprintf("This field must not exceed %s characters.", param)
	case "len":
		return fmt.Sprintf("This field must have exactly %s characters.", param)
	default:
		return fmt.Sprintf("Invalid value for this field: %s", tag)
	}
}
