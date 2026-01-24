package entity

import (
	"api/internal/entity/enum"
	"time"

	"gorm.io/gorm"
)

type EmployeeAttendance struct {
	ID     int                   `gorm:"primaryKey"`
	Date   time.Time             `gorm:"type:date;not null"`
	Status enum.AttendanceStatus `gorm:"type:absen_status;column:status;not null"`

	EmployeeId int       `gorm:"column:employee_id;not null"`
	Employee   *Employee `gorm:"foreignKey:EmployeeId;references:ID"`

	PeriodId int     `gorm:"column:period_id;not null"`
	Period   *Period `gorm:"foreignKey:PeriodId;references:ID"`

	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoCreateTime:milli;autoUpdateTime:milli"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (ea *EmployeeAttendance) TableName() string {
	return "employee_attendances"
}
