package model

import (
	"api/internal/entity/enum"
	"time"

	"github.com/google/uuid"
)

type FindAllUserRequest struct {
	Name     string          `json:"name" validate:"omitempty,max=100"`
	Username string          `json:"username" validate:"omitempty,max=100"`
	Phone    string          `json:"phone" validate:"omitempty"`
	Roles    []enum.UserRole `json:"roles" validate:"omitempty,userrole"`
	Page     int             `json:"page"`
	PerPage  int             `json:"perPage" validate:"max=100"`
}

type RegisterUserRequest struct {
	Name     string        `json:"name" validate:"required,max=100"`
	Username string        `json:"username" validate:"required"`
	Password string        `json:"password" validate:"required,min=6"`
	Phone    string        `json:"phone" validate:"required,numeric,max=15"`
	Role     enum.UserRole `json:"role" validate:"required,userrole"`
}

type UpdateUserRequest struct {
	ID       uuid.UUID     `json:"id" validate:"required,uuid"`
	Name     string        `json:"name" validate:"omitempty,required,max=100"`
	Username string        `json:"username" validate:"omitempty,required"`
	Password string        `json:"password" validate:"omitempty,min=6"`
	Phone    string        `json:"phone" validate:"omitempty,required,numeric,max=15"`
	Role     enum.UserRole `json:"role" validate:"omitempty,required,userrole"`
}

type DeleteUserRequest struct {
	ID uuid.UUID `json:"id" validate:"required,uuid"`
}

type LoginUserRequest struct {
	Password string `json:"password" validate:"required,min=6"`
	Username string `json:"username" validate:"required"`
}

type VerifyUserRequest struct {
	Token string `validate:"required,max=200"`
}

type UserResponse struct {
	ID        uuid.UUID     `json:"id,omitempty"`
	Name      string        `json:"name"`
	Username  string        `json:"username"`
	Phone     string        `json:"phone"`
	Role      enum.UserRole `json:"role"`
	CreatedAt time.Time     `json:"createdAt"`
}
