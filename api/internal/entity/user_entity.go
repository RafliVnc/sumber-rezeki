package entity

import (
	"api/internal/entity/enum"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey"`
	Name      string         `gorm:"column:name"`
	Username  string         `gorm:"column:username"`
	Phone     string         `gorm:"column:phone"`
	Role      enum.UserRole  `gorm:"type:user_role;column:role;not null"`
	Password  string         `gorm:"column:password"`
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoCreateTime:milli;autoUpdateTime:milli"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (u *User) TableName() string {
	return "users"
}
