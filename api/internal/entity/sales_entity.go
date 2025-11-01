package entity

import (
	"time"

	"gorm.io/gorm"
)

type Sales struct {
	ID    int    `gorm:"primaryKey;autoIncrement"`
	Name  string `gorm:"column:name"`
	Phone string `gorm:"column:phone"`

	Routes []Route `gorm:"many2many:sales_routes;"`

	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoCreateTime:milli;autoUpdateTime:milli"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (s *Sales) TableName() string {
	return "sales"
}
