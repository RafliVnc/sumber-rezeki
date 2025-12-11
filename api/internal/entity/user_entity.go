package entity

import (
	"api/internal/entity/enum"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey"`
	Name      string         `gorm:"column:name;not null"`
	Username  string         `gorm:"column:username;not null"`
	Phone     string         `gorm:"column:phone;not null"`
	Role      enum.UserRole  `gorm:"type:UserRole;column:role;not null"`
	Password  string         `gorm:"column:password;not null"`
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoCreateTime:milli;autoUpdateTime:milli"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`

	Periods        []Period        `gorm:"foreignKey:ClosedBy;references:ID"`
	PeriodClosures []PeriodClosure `gorm:"foreignKey:ClosedBy;references:ID"`
}

func (u *User) TableName() string {
	return "users"
}
