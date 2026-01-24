package entity

import (
	"time"

	"gorm.io/gorm"
)

type Sales struct {
	ID    int    `gorm:"primaryKey;autoIncrement"`
	Phone string `gorm:"column:phone"`

	EmployeeId int      `gorm:"unique"`
	Employee   Employee `gorm:"foreignKey:EmployeeId"`

	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoCreateTime:milli;autoUpdateTime:milli"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`

	Routes []Route `gorm:"many2many:sales_routes;"`
}

func (s *Sales) TableName() string {
	return "sales"
}
