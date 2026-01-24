package model

import "api/internal/entity/enum"

type EmployeeResponse struct {
	ID           int                          `json:"id,omitempty"`
	Name         string                       `json:"name"`
	Salary       float64                      `json:"salary"`
	SupervisorId *int                         `json:"supervisorId,omitempty"`
	Role         string                       `json:"role"`
	Sales        *SalesResponse               `json:"Sales,omitempty"`
	Attendaces   []EmployeeAttendanceResponse `json:"Attendaces,omitempty"`
}

type FindAllEmployeeRequest struct {
	Page    int     `json:"page" validate:"omitempty,max=100"`
	PerPage int     `json:"perPage" validate:"omitempty"`
	Name    string  `json:"name" validate:"omitempty,max=100"`
	Salary  float64 `json:"salary" validate:"omitempty"`
	// TODO: Create EmployeeRole validation
	Roles []enum.EmployeeRole `json:"roles" validate:"omitempty"`
}

type CreateEmployeeRequest struct {
	Name         string            `json:"name" validate:"required,max=100"`
	Salary       float64           `json:"salary" validate:"required"`
	Role         enum.EmployeeRole `json:"role" validate:"required,oneof='WAREHOUSE_HEAD' 'SALES' 'DRIVER' 'HELPER' 'TREASURER' 'STAFF'"`
	SupervisorId int               `json:"supervisorId" validate:"required_if=Role HELPER,required_if=Role DRIVER"`
	Phone        string            `json:"phone" validate:"required_if=Role SALES"`
	RouteIDs     []int             `json:"routeIds" validate:"required_if=Role SALES"`
}

type UpdateEmployeeRequest struct {
	ID           int               `json:"id" validate:"required,gt=0"`
	Name         string            `json:"name" validate:"required,max=100"`
	Salary       float64           `json:"salary" validate:"required"`
	Role         enum.EmployeeRole `json:"role" validate:"required,oneof='WAREHOUSE_HEAD' 'SALES' 'DRIVER' 'HELPER' 'TREASURER' 'STAFF'"`
	SupervisorId int               `json:"supervisorId" validate:"required_if=Role HELPER,required_if=Role DRIVER"`
	Phone        string            `json:"phone" validate:"required_if=Role SALES"`
	RouteIDs     *[]int            `json:"routeIds" validate:"required_if=Role SALES"`
}

type DeleteEmployeeRequest struct {
	ID int `json:"id" validate:"required,gt=0"`
}

type FindByIdEmployeeRequest struct {
	ID int `json:"id" validate:"required,gt=0"`
}

type FindAllEmployeeWithAttendanceRequest struct {
	StartDate string `json:"startDate" validate:"required"`
	EndDate   string `json:"endDate" validate:"required"`
}
