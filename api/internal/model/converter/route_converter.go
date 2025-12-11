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

func ToRouteResponseWithSalesCount(route *entity.Route, salesCount int) *model.RouteResponseWithSalesCount {
	return &model.RouteResponseWithSalesCount{
		ID:          route.ID,
		Name:        route.Name,
		Description: route.Description,
		SalesCount:  salesCount,
	}
}
