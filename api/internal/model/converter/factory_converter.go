package converter

import (
	"api/internal/entity"
	"api/internal/model"
)

func ToFactoryResponse(factory *entity.Factory) *model.FactoryResponse {
	return &model.FactoryResponse{
		ID:          factory.ID,
		Name:        factory.Name,
		Phone:       factory.Phone,
		DueDate:     factory.DueDate,
		Description: factory.Description,
	}
}
