package model

import "api/internal/entity/enum"

type FindAllVehicleRequest struct {
	Search  string             `json:"search" validate:"omitempty,max=100"`
	Types   []enum.VehicleType `json:"type" validate:"omitempty,dive,oneof='PICKUP' 'TRONTON' 'TRUCK'"`
	Page    int                `json:"page"`
	PerPage int                `json:"perPage" validate:"max=100"`
}

type CreateVehicleRequest struct {
	Plate string           `json:"plate" validate:"required,max=20"`
	Type  enum.VehicleType `json:"type" validate:"required,oneof='PICKUP' 'TRONTON' 'TRUCK'"`
}

type UpdateVehicleRequest struct {
	ID    int64            `json:"id" validate:"required,gt=0"`
	Plate string           `json:"plate" validate:"required,max=20"`
	Type  enum.VehicleType `json:"type" validate:"required,oneof='PICKUP' 'TRONTON' 'TRUCK'"`
}

type DeleteVehicleRequest struct {
	ID int64 `json:"id" validate:"required,gt=0"`
}

type VehicleResponse struct {
	ID    int64            `json:"id"`
	Plate string           `json:"plate"`
	Type  enum.VehicleType `json:"type"`
}
