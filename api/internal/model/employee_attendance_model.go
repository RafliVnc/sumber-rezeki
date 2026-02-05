package model

import "time"

type EmployeeAttendanceResponse struct {
	ID         int       `json:"id"`
	Status     string    `json:"status"`
	Date       time.Time `json:"date"`
	EmployeeId int       `json:"employeeId,omitempty"`
	PeriodId   int       `json:"periodId"`
}

type UpsertEmployeeAttendanceRequest struct {
	Attendances []AttendanceAction `json:"attendances" validate:"required,dive"`
}

type AttendanceAction struct {
	Action    string               `json:"action" validate:"required,oneof=update delete"`
	Date      string               `json:"date" validate:"required"`
	Employees []EmployeeAttendance `json:"employees" validate:"omitempty,dive"`
}

type EmployeeAttendance struct {
	ID     int    `json:"id" validate:"required"`
	Status string `json:"status" validate:"required,oneof=PRESENT ABSENT LEAVE SICK"`
}

type FindAllEmployeeAttendanceRequest struct {
	StartDate string `json:"startDate" validate:"required"`
	EndDate   string `json:"endDate" validate:"required"`
}
