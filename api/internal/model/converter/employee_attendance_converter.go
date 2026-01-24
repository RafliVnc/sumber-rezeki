package converter

import (
	"api/internal/entity"
	"api/internal/model"
)

func ToEmployeeAttendanceResponse(EmployeeAttendance *entity.EmployeeAttendance) *model.EmployeeAttendanceResponse {
	return &model.EmployeeAttendanceResponse{
		ID:       EmployeeAttendance.ID,
		Status:   string(EmployeeAttendance.Status),
		Date:     EmployeeAttendance.Date,
		PeriodId: EmployeeAttendance.PeriodId,
	}
}
