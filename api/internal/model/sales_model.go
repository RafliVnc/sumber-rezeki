package model

import (
	"time"
)

type FindAllSalesRequest struct {
	Search   string `json:"search" validate:"omitempty,max=100"`
	RouteIDs []int  `json:"routeIds" validate:"omitempty,min=1,dive,gt=0"`
	Page     int    `json:"page"`
	PerPage  int    `json:"perPage" validate:"max=100"`
}

type CreateSalesRequest struct {
	Phone      string `json:"phone" validate:"required,numeric,max=15"`
	EmployeeID int    `json:"employeeId"  validate:"required"`
	RouteIDs   []int  `json:"routeIds" validate:"omitempty,dive,min=1"`
}

type UpdateSalesRequest struct {
	ID         int    `json:"id" validate:"required,gt=0"`
	Phone      string `json:"phone" validate:"omitempty,max=20"`
	EmployeeID int    `json:"employeeId"  validate:"omitempty"`
	RouteIDs   []int  `json:"routeIds" validate:"omitempty,min=1,dive,gt=0"`
}

type DeleteSalesRequest struct {
	ID int `json:"id" validate:"required,gt=0"`
}

type SalesResponse struct {
	ID        int       `json:"id,omitempty"`
	Phone     string    `json:"phone"`
	CreatedAt time.Time `json:"createdAt"`

	EmployeeID *int              `json:"employeeId,omitempty"`
	Employee   *EmployeeResponse `json:"Employee,omitempty"`
	Routes     []RouteResponse   `json:"Routes"`
}
