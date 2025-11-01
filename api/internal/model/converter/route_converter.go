package converter

import (
	"api/internal/entity"
	"api/internal/model"
)

func ToRouteResponse(route *entity.Route) *model.RouteResponse {
	return &model.RouteResponse{
		ID:          route.ID,
		Name:        route.Name,
		Description: route.Description,
	}
}
