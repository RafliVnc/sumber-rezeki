package entity

import (
	"time"

	"gorm.io/gorm"
)

type Route struct {
	ID          int    `gorm:"primaryKey;autoIncrement"`
	Name        string `gorm:"column:name"`
	Description string `gorm:"column:description"`

	Sales []Sales `gorm:"many2many:sales_routes;"`

	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoCreateTime:milli;autoUpdateTime:milli"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (r *Route) TableName() string {
	return "routes"
}
