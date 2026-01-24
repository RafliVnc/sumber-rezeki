package model

import "time"

type EmployeeAttendanceResponse struct {
	ID         int       `json:"id"`
	Status     string    `json:"status"`
	Date       time.Time `json:"date"`
	EmployeeId int       `json:"employeeId,omitempty"`
	PeriodId   int       `json:"periodId"`
}

type CreateEmployeeAttendanceRequest struct {
	Action     string `json:"action" validate:"required,oneof='upsert' 'delete'"`
	Date       string `json:"date" validate:"required"`
	Status     string `json:"status" validate:"required,oneof= PRESENT ABSENT LEAVE SICK"`
	EmployeeId int    `jason:"employeeId" validate:"required"`
}

type FindAllEmployeeAttendanceRequest struct {
	StartDate string `json:"startDate" validate:"required"`
	EndDate   string `json:"endDate" validate:"required"`
}
