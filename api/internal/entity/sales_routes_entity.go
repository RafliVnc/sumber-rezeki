package entity

import (
	"time"

	"gorm.io/gorm"
)

type SalesRoutes struct {
	ID      int `gorm:"type:number;primaryKey"`
	SalesId int `gorm:"column:sales_id;type:number"`
	RouteId int `gorm:"column:route_id;type:number"`

	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoCreateTime:milli;autoUpdateTime:milli"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (sr *SalesRoutes) TableName() string {
	return "sales_routes"
}
