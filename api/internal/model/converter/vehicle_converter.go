package converter

import (
	"api/internal/entity"
	"api/internal/model"
)

func ToVehicleResponse(vehicle *entity.Vehicle) *model.VehicleResponse {
	return &model.VehicleResponse{
		ID:    vehicle.ID,
		Plate: vehicle.Plate,
		Type:  vehicle.Type,
	}
}
