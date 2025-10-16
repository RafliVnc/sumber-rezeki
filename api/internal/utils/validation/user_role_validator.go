package validation

import (
	"api/internal/entity/enum"

	"github.com/go-playground/validator/v10"
)

// validate user role
func RegisterUserRoleValidation(validate *validator.Validate) {
	validate.RegisterValidation("userrole", func(fl validator.FieldLevel) bool {
		role, ok := fl.Field().Interface().(enum.UserRole)
		if !ok {
			return false
		}

		switch role {
		case enum.Owner, enum.SuperAdmin, enum.Treasurer, enum.WarehouseHead:
			return true
		default:
			return false
		}
	})
}
