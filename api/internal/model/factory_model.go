package model

type FindAllFactoryRequest struct {
	Search  string `json:"search" validate:"omitempty,max=100"`
	Page    int    `json:"page"`
	PerPage int    `json:"perPage" validate:"max=100"`
}

type CreateFactoryRequest struct {
	Name        string  `json:"name" validate:"required,max=100"`
	Phone       string  `json:"phone" validate:"required,max=20"`
	DueDate     int64   `json:"dueDate" validate:"required"`
	Description *string `json:"description,omitempty"`
}

type UpdateFactoryRequest struct {
	ID          int64   `json:"id" validate:"required,gt=0"`
	Name        string  `json:"name" validate:"required,max=100"`
	Phone       string  `json:"phone" validate:"required,max=20"`
	DueDate     int64   `json:"dueDate" validate:"required"`
	Description *string `json:"description,omitempty"`
}

type DeleteFactoryRequest struct {
	ID int64 `json:"id" validate:"required,gt=0"`
}

type FactoryResponse struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Phone       string  `json:"phone"`
	DueDate     int64   `json:"dueDate"`
	Description *string `json:"description"`
}
