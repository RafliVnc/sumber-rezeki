package test

import (
	"api/internal/model"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateRoute(t *testing.T) {
	defer ClearAll()

	token, err := GenerateTokenHelper()
	assert.Nil(t, err)

	requestBody := model.CreateRouteRequest{
		Name:        "Route Test 1",
		Description: "Description for route test 1",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPost, "/api/routes", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", token)
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.RouteResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusCreated, response.StatusCode)
	assert.NotNil(t, responseBody.Data.ID)
	assert.Equal(t, requestBody.Name, responseBody.Data.Name)
	assert.Equal(t, requestBody.Description, responseBody.Data.Description)
}

func TestCreateRouteDuplicateName(t *testing.T) {
	defer ClearAll()

	token, err := GenerateTokenHelper()
	assert.Nil(t, err)

	CreateRoutes(1)

	// Create duplicate route
	requestBody := model.CreateRouteRequest{
		Name:        "Created Route 0",
		Description: "Description for route test 0",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPost, "/api/routes", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", token)
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ErrorResponse)
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Equal(t, http.StatusBadRequest, responseBody.Code)
	assert.NotEmpty(t, responseBody.Message)
}

func TestGetAllRoutes(t *testing.T) {
	defer ClearAll()

	token, err := GenerateTokenHelper()
	assert.Nil(t, err)

	// Create multiple routes
	CreateRoutes(15)

	request := httptest.NewRequest(http.MethodGet, "/api/routes?page=1&perPage=10", nil)
	request.Header.Set("Authorization", token)
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[[]model.RouteResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.NotNil(t, responseBody.Data)
	assert.LessOrEqual(t, len(responseBody.Data), 10)
	assert.NotNil(t, responseBody.Paging)
	assert.Equal(t, 1, responseBody.Paging.Page)
	assert.Equal(t, 10, responseBody.Paging.PerPage)
}

func TestGetAllRoutesSearch(t *testing.T) {
	defer ClearAll()

	token, err := GenerateTokenHelper()
	assert.Nil(t, err)

	CreateRoutes(15)

	request := httptest.NewRequest(http.MethodGet, "/api/routes?page=1&perPage=10&search=Created+Route+8", nil)
	request.Header.Set("Authorization", token)
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[[]model.RouteResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.NotNil(t, responseBody.Data)
	assert.GreaterOrEqual(t, len(responseBody.Data), 1)
	assert.NotNil(t, responseBody.Paging)
	assert.Equal(t, 1, responseBody.Paging.Page)
	assert.Equal(t, 10, responseBody.Paging.PerPage)
}

func TestUpdateRoute(t *testing.T) {
	defer ClearAll()

	token, err := GenerateTokenHelper()
	assert.Nil(t, err)

	// Create route first
	routes := CreateRoutes(1)
	routeID := routes[0].ID

	requestBody := model.UpdateRouteRequest{
		ID:          routeID,
		Name:        "Updated Route Name",
		Description: "Updated description",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/routes/%d", routeID), strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", token)
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.RouteResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, routeID, responseBody.Data.ID)
	assert.Equal(t, requestBody.Name, responseBody.Data.Name)
	assert.Equal(t, requestBody.Description, responseBody.Data.Description)
}

func TestUpdateRouteNotFound(t *testing.T) {
	defer ClearAll()

	token, err := GenerateTokenHelper()
	assert.Nil(t, err)

	requestBody := model.UpdateRouteRequest{
		ID:          99999,
		Name:        "Non Existent Route",
		Description: "This should fail",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPut, "/api/routes/99999", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", token)
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ErrorResponse)
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusNotFound, response.StatusCode)
	assert.Equal(t, http.StatusNotFound, responseBody.Code)
	assert.NotEmpty(t, responseBody.Message)
}

func TestUpdateRouteDuplicateName(t *testing.T) {
	defer ClearAll()

	token, err := GenerateTokenHelper()
	assert.Nil(t, err)

	// Create two routes
	routes := CreateRoutes(2)

	// Try to update second route with first route's name
	requestBody := model.UpdateRouteRequest{
		ID:          routes[1].ID,
		Name:        routes[0].Name,
		Description: "Updated description",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/routes/%d", routes[1].ID), strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", token)
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ErrorResponse)
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Equal(t, http.StatusBadRequest, responseBody.Code)
	assert.NotEmpty(t, responseBody.Message)
}

func TestDeleteRoute(t *testing.T) {
	defer ClearAll()

	token, err := GenerateTokenHelper()
	assert.Nil(t, err)

	// Create route first
	routes := CreateRoutes(1)
	routeID := routes[0].ID

	request := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/routes/%d", routeID), nil)
	request.Header.Set("Authorization", token)
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
}

func TestDeleteRouteNotFound(t *testing.T) {
	defer ClearAll()

	token, err := GenerateTokenHelper()
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodDelete, "/api/routes/99999", nil)
	request.Header.Set("Authorization", token)
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ErrorResponse)
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusNotFound, response.StatusCode)
	assert.Equal(t, http.StatusNotFound, responseBody.Code)
	assert.NotEmpty(t, responseBody.Message)
}

func TestDeleteRouteWithSales(t *testing.T) {
	defer ClearAll()

	token, err := GenerateTokenHelper()
	assert.Nil(t, err)

	// Create route and assign to sales
	routes := CreateRoutes(1)
	routeID := routes[0].ID

	// Create sales with this route
	CreateSalesWithRoutes(1, []int{routeID})

	request := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/routes/%d", routeID), nil)
	request.Header.Set("Authorization", token)
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ErrorResponse)
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Equal(t, http.StatusBadRequest, responseBody.Code)
	assert.NotEmpty(t, responseBody.Message)
	assert.Contains(t, responseBody.Message, "sales")
}
