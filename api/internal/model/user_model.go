package model

import (
	"api/internal/entity/enum"

	"github.com/google/uuid"
)

type RegisterUserRequest struct {
	Name     string        `json:"name" validate:"required,max=100"`
	Username string        `json:"username" validate:"required"`
	Password string        `json:"password" validate:"required,min=6"`
	Phone    string        `json:"phone" validate:"required,numeric,max=15"`
	Role     enum.UserRole `json:"role" validate:"required,userrole"`
}

type UpdateUserRequest struct {
	ID       uuid.UUID `json:"id" validate:"required,gt=0"`
	Name     string    `json:"name,omitempty" validate:"max=100"`
	Username string    `json:"username,omitempty" validate:"max=100"`
}

type LoginUserRequest struct {
	Password string `json:"password" validate:"required,min=6"`
	Username string `json:"username" validate:"required"`
}

type DeleteUserRequest struct {
	ID uuid.UUID `json:"id" validate:"required,gt=0"`
}

type UserResponse struct {
	ID       uuid.UUID     `json:"id,omitempty"`
	Name     string        `json:"name,omitempty"`
	Username string        `json:"username,omitempty"`
	Phone    string        `json:"phone,omitempty"`
	Role     enum.UserRole `json:"role,omitempty"`
}

type LoginUserResponse struct {
	User  UserResponse
	Token string `json:"token"`
}

type FindAllUserRequest struct {
	Name     string `json:"name" validate:"max=100"`
	Username string `json:"username" validate:"max=100"`
	Page     int    `json:"page"`
	PerPage  int    `json:"per_page" validate:"max=100"`
}

type VerifyUserRequest struct {
	Token string `validate:"required,max=200"`
}

type VerifyTokenRequest struct {
	Token string `validate:"required,max=100"`
}
