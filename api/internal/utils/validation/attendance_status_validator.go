package validation

import (
	"api/internal/entity/enum"
	"reflect"

	"github.com/go-playground/validator/v10"
)

// validate absent status (single or slice)
func RegisterAbsenStatusValidation(validate *validator.Validate) {
	validate.RegisterValidation("AttendanceStatus", func(fl validator.FieldLevel) bool {
		field := fl.Field()

		// Handle slice of AttendanceStatus
		if field.Kind() == reflect.Slice {
			for i := 0; i < field.Len(); i++ {
				status, ok := field.Index(i).Interface().(enum.AttendanceStatus)
				if !ok {
					return false
				}

				if !isValidStatus(status) {
					return false
				}
			}
			return true
		}

		// Handle single AttendanceStatus
		status, ok := field.Interface().(enum.AttendanceStatus)
		if !ok {
			return false
		}

		return isValidStatus(status)
	})
}

// Helper function to check valid status
func isValidStatus(status enum.AttendanceStatus) bool {
	switch status {
	case enum.PERMIT, enum.PRESENT:
		return true
	default:
		return false
	}
}
