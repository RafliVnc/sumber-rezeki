package converter

import (
	"api/internal/entity"
	"api/internal/model"
)

func ToSalesResponse(sales *entity.Sales) *model.SalesResponse {
	routes := make([]model.RouteResponse, len(sales.Routes))
	for i, r := range sales.Routes {
		routes[i] = *ToRouteResponse(&r)
	}

	return &model.SalesResponse{
		ID:        sales.ID,
		Name:      sales.Name,
		Phone:     sales.Phone,
		CreatedAt: sales.CreatedAt,
		Routes:    routes,
	}
}
