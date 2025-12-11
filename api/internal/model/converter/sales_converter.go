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
		EmployeeID: &sales.EmployeeID,
		Employee:   &model.EmployeeResponse{Name: sales.Employee.Name},
	}
}
