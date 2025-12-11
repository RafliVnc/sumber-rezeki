package converter

import (
	"api/internal/entity"
	"api/internal/model"
)

func ToEmployeeResponse(employee *entity.Employee) *model.EmployeeResponse {
	return &model.EmployeeResponse{
		ID:           employee.ID,
		Name:         employee.Name,
		Salary:       employee.Salary,
		Role:         string(employee.Role),
		SupervisorId: employee.SupervisorId,
	}
}
