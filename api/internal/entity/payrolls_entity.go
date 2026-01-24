package entity

import (
	"api/internal/entity/enum"
	"time"

	"gorm.io/gorm"
)

type Payroll struct {
	ID             int                `gorm:"primaryKey;autoIncrement"`
	BaseSalary     float64            `gorm:"column:base_salary;not null"`
	AttendanceDays int                `gorm:"column:attendance_days;not null"`
	Deductions     float64            `gorm:"column:deductions;not null:default:0"`
	Bonuses        float64            `gorm:"column:bonuses;not null:default:0"`
	ModuleType     enum.PayrollModule `gorm:"column:module_type;not null"`
	Notes          string             `gorm:"column:notes"`
	IsPaid         bool               `gorm:"column:is_paid;not null;default:false"`
	PaidAt         *time.Time         `gorm:"column:paid_at"`

	EmployeeId int       `gorm:"column:employee_id;not null"`
	Employee   *Employee `gorm:"foreignKey:EmployeeId;references:ID"`
	PeriodID   int       `gorm:"column:period_id;not null"`
	Period     *Period   `gorm:"foreignKey:PeriodID;references:ID"`

	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoCreateTime:milli;autoUpdateTime:milli"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (p *Payroll) TableName() string {
	return "payrolls"
}
