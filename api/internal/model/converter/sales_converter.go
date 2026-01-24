package converter

import (
	"api/internal/entity"
	"api/internal/model"
)

func ToSalesResponse(sales *entity.Sales) *model.SalesResponse {

	return &model.SalesResponse{
		ID:         sales.ID,
		Phone:      sales.Phone,
		CreatedAt:  sales.CreatedAt,
		EmployeeId: &sales.EmployeeId,
		Employee:   &model.EmployeeResponse{Name: sales.Employee.Name},
	}
}
