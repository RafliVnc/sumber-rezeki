package entity

import (
	"time"

	"gorm.io/gorm"
)

type Factory struct {
	ID          int64   `gorm:"primaryKey;autoIncrement;column:id"`
	Name        string  `gorm:"column:name;type:varchar(100);not null"`
	DueDate     int64   `gorm:"column:due_date;not null"`
	Phone       string  `gorm:"column:phone;type:varchar(20);not null"`
	Description *string `gorm:"column:description;type:text"`

	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoCreateTime:milli;autoUpdateTime:milli"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (f *Factory) TableName() string {
	return "factories"
}
