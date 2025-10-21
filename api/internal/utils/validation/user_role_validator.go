package validation

import (
	"api/internal/entity/enum"
	"reflect"

	"github.com/go-playground/validator/v10"
)

// validate user role (single or slice)
func RegisterUserRoleValidation(validate *validator.Validate) {
	validate.RegisterValidation("userrole", func(fl validator.FieldLevel) bool {
		field := fl.Field()

		// Handle slice of UserRole
		if field.Kind() == reflect.Slice {
			for i := 0; i < field.Len(); i++ {
				role, ok := field.Index(i).Interface().(enum.UserRole)
				if !ok {
					return false
				}

				if !isValidRole(role) {
					return false
				}
			}
			return true
		}

		// Handle single UserRole
		role, ok := field.Interface().(enum.UserRole)
		if !ok {
			return false
		}

		return isValidRole(role)
	})
}

// Helper function to check valid role
func isValidRole(role enum.UserRole) bool {
	switch role {
	case enum.OWNER, enum.SUPER_ADMIN, enum.TREASURER, enum.WAREHOUSE_HEAD:
		return true
	default:
		return false
	}
}
