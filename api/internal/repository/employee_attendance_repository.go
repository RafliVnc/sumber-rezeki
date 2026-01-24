package repository

import (
	"github.com/sirupsen/logrus"
)

type EmployeeAttendanceRepository interface {
}

type employeeAttendanceRepositoryImpl struct {
	Log *logrus.Logger
}

func NewEmployeeAttendanceRepository(log *logrus.Logger) EmployeeAttendanceRepository {
	return &employeeAttendanceRepositoryImpl{
		Log: log,
	}
}
