package entity

import (
	"api/internal/entity/enum"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Period struct {
	ID         int             `gorm:"primaryKey;autoIncrement"`
	Type       enum.PeriodType `gorm:"column:type;not null"`
	StartDate  time.Time       `gorm:"column:start_date;not null"`
	EndDate    time.Time       `gorm:"column:end_date;not null"`
	WeekNumber int             `gorm:"column:week_number;not null"`
	Month      int             `gorm:"column:month;not null"`
	Year       int             `gorm:"column:year;not null"`
	IsActive   bool            `gorm:"column:is_active;not null;default:true"`

	IsClosed     bool       `gorm:"column:is_closed;not null;default:false"`
	ClosedBy     *uuid.UUID `gorm:"type:uuid;column:closed_by"`
	ClosedAt     *time.Time `gorm:"column:closed_at"`
	ClosedByUser *User      `gorm:"foreignKey:ClosedBy;references:ID"`

	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoCreateTime:milli;autoUpdateTime:milli"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (p *Period) TableName() string {
	return "periods"
}
