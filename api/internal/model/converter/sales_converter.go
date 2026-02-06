package converter

import (
	"api/internal/entity"
	"api/internal/model"
)

func ToSalesResponse(sales *entity.Sales) *model.SalesResponse {

	return &model.SalesResponse{
		ID:        sales.ID,
		Phone:     sales.Phone,
		CreatedAt: sales.CreatedAt,
		Employee:  model.EmployeeResponse{ID: sales.Employee.ID, Name: sales.Employee.Name},
	}
}
