package entity

import (
	"api/internal/entity/enum"
	"time"

	"gorm.io/gorm"
)

type Vehicle struct {
	ID    int64            `gorm:"primaryKey;autoIncrement;column:id"`
	Plate string           `gorm:"column:plate;type:varchar(100);not null"`
	Type  enum.VehicleType `gorm:"column:type;type:VehicleType;not null"`

	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime:milli"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`

	Histories []VehicleHistory `gorm:"foreignKey:VehicleID"`
}

func (v *Vehicle) TableName() string {
	return "vehicles"
}
