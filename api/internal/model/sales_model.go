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

type UpdateSalesRequest struct {
	ID       int    `json:"id" validate:"required,gt=0"`
	Name     string `json:"name" validate:"required,max=100"`
	Phone    string `json:"phone" validate:"required,max=20"`
	RouteIDs []int  `json:"routeIds" validate:"omitempty,min=1,dive,gt=0"`
}

type DeleteSalesRequest struct {
	ID int `json:"id" validate:"required,gt=0"`
}

type SalesResponse struct {
	ID        int       `json:"id,omitempty"`
	Phone     string    `json:"phone"`
	CreatedAt time.Time `json:"createdAt"`

	Employee EmployeeResponse `json:"Employee"`
	Routes   []RouteResponse  `json:"Routes,omitempty"`
}
