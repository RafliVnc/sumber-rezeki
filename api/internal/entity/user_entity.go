package entity

import (
	"api/internal/entity/enum"

	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID     `gorm:"type:uuid;primaryKey"`
	Name     string        `gorm:"column:name"`
	Username string        `gorm:"column:username"`
	Phone    string        `gorm:"column:phone"`
	Role     enum.UserRole `gorm:"type:user_role;column:role;not null"`
	Password string        `gorm:"column:password"`
}

func (u *User) TableName() string {
	return "users"
}
