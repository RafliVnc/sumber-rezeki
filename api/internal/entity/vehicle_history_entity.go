package entity

import (
	"api/internal/entity/enum"
	"time"

	"gorm.io/gorm"
)

type VehicleHistory struct {
	ID          int64                   `gorm:"primaryKey;autoIncrement;column:id"`
	Date        time.Time               `gorm:"column:date;not null"`
	Description string                  `gorm:"column:description;type:text;not null"`
	Type        enum.VehicleHistoryType `gorm:"column:type;type:VehicleHistoryType;not null"`

	Amount float64 `gorm:"column:amount;type:decimal(12,2);not null;default:0"`
	Profit *int    `gorm:"column:profit"`
	Sack   *int    `gorm:"column:sack"`

	VehicleID int64   `gorm:"column:vehicle_id;not null"`
	Vehicle   Vehicle `gorm:"foreignKey:VehicleID;references:ID;constraint:OnDelete:RESTRICT"`

	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime:milli"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (vh *VehicleHistory) TableName() string {
	return "vehicle_history"
}
