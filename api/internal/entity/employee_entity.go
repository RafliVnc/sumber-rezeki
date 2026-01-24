package entity

import (
	"api/internal/entity/enum"
	"time"

	"gorm.io/gorm"
)

type Employee struct {
	ID       int               `gorm:"primaryKey;autoIncrement"`
	Name     string            `gorm:"column:name;not null"`
	Salary   float64           `gorm:"column:salary;not null"`
	Role     enum.EmployeeRole `gorm:"column:role;not null"`
	JoinDate time.Time         `gorm:"column:join_date;not null"`

	SupervisorId *int       `gorm:"column:supervisor_id"`
	Supervisor   *Employee  `gorm:"foreignKey:SupervisorId;references:ID"`
	Subordinates []Employee `gorm:"foreignKey:SupervisorId;references:ID"`

	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoCreateTime:milli;autoUpdateTime:milli"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`

	EmployeeAttendance []EmployeeAttendance `gorm:"foreignKey:EmployeeId;references:ID"`
	Sales              *Sales               `gorm:"foreignKey:EmployeeId;references:ID"`
	Payrolls           []Payroll            `gorm:"foreignKey:EmployeeId;references:ID"`
}

func (e *Employee) TableName() string {
	return "employees"
}
