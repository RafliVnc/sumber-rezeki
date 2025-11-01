package model

type RouteResponse struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

type CreateRouteRequest struct {
	Name        string `json:"name" validate:"required,max=100"`
	Description string `json:"description,omitempty"`
}

type UpdateRouteRequest struct {
	ID          int    `json:"id" validate:"required,gt=0"`
	Name        string `json:"name" validate:"omitempty,max=100"`
	Description string `json:"description,omitempty"`
}

type DeleteRouteRequest struct {
	ID int `json:"id" validate:"required,gt=0"`
}

type FindAllRouteRequest struct {
	Search  string `json:"search" validate:"omitempty,max=100"`
	Page    int    `json:"page"`
	PerPage int    `json:"perPage" validate:"max=100"`
}
