package model

type RegisterUserRequest struct {
	Name     string `json:"name" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}
type UpdateUserRequest struct {
	ID    int    `json:"id" validate:"required,gt=0"`
	Name  string `json:"name,omitempty" validate:"max=100"`
	Email string `json:"email,omitempty" validate:"email"`
}

type DeleteUserRequest struct {
	ID int `json:"id" validate:"required,gt=0"`
}

type UserResponse struct {
	ID    int    `json:"id,omitempty"`
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
}

type FindAllUserRequest struct {
	Name    string `json:"name" validate:"max=100"`
	Email   string `json:"email" validate:"max=100"`
	Page    int    `json:"page"`
	PerPage int    `json:"per_page" validate:"max=100"`
}

type VerifyUserRequest struct {
	Token string `validate:"required,max=100"`
}
