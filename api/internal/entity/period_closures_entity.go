package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PeriodClosure struct {
	ID         int    `gorm:"primaryKey;autoIncrement"`
	ModuleName string `gorm:"column:module_name;not null"`
	Notes      string `gorm:"column:notes;not null"`
	PrintCount int    `gorm:"column:print_count;not null;default:0"`

	PeriodID int    `gorm:"column:period_id;not null"`
	Period   Period `gorm:"foreignKey:PeriodID;references:ID"`

	IsClosed     bool       `gorm:"column:is_closed;not null;default:false"`
	ClosedBy     *uuid.UUID `gorm:"type:uuid;column:closed_by"`
	ClosedAt     *time.Time `gorm:"column:closed_at"`
	ClosedByUser *User      `gorm:"foreignKey:ClosedBy;references:ID"`

	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoCreateTime:milli;autoUpdateTime:milli"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (pc *PeriodClosure) TableName() string {
	return "period_closures"
}
