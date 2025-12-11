package test

import (
	"api/internal/model"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// func TestCreateSales(t *testing.T) {
// 	defer ClearAll()

// 	token, err := GenerateTokenHelper()
// 	assert.Nil(t, err)

// 	requestBody := model.CreateSalesRequest{
// 		Name:  "Created Sales Test",
// 		Phone: "08123456789",
// 	}

// 	bodyJson, err := json.Marshal(requestBody)
// 	assert.Nil(t, err)

// 	request := httptest.NewRequest(http.MethodPost, "/api/sales", strings.NewReader(string(bodyJson)))
// 	request.Header.Set("Content-Type", "application/json")
// 	request.Header.Set("Authorization", token)
// 	request.Header.Set("Accept", "application/json")

// 	response, err := app.Test(request)
// 	assert.Nil(t, err)

// 	bytes, err := io.ReadAll(response.Body)
// 	assert.Nil(t, err)

// 	responseBody := new(model.WebResponse[model.SalesResponse])
// 	err = json.Unmarshal(bytes, responseBody)
// 	assert.Nil(t, err)

// 	assert.Equal(t, http.StatusCreated, response.StatusCode)
// 	assert.NotNil(t, responseBody.Data.ID)
// 	assert.Equal(t, requestBody.Name, responseBody.Data.Name)
// 	assert.Equal(t, requestBody.Phone, responseBody.Data.Phone)
// 	assert.NotEmpty(t, responseBody.Data.CreatedAt)
// }

// func TestCreateSalesDuplicatePhone(t *testing.T) {
// 	defer ClearAll()

// 	token, err := GenerateTokenHelper()
// 	assert.Nil(t, err)

// 	CreateSales(1)

// 	requestBody := model.CreateSalesRequest{
// 		Name:  "Created Sales Test",
// 		Phone: "1234560",
// 	}

// 	bodyJson, err := json.Marshal(requestBody)
// 	assert.Nil(t, err)

// 	request := httptest.NewRequest(http.MethodPost, "/api/sales", strings.NewReader(string(bodyJson)))
// 	request.Header.Set("Content-Type", "application/json")
// 	request.Header.Set("Authorization", token)
// 	request.Header.Set("Accept", "application/json")

// 	response, err := app.Test(request)
// 	assert.Nil(t, err)

// 	bytes, err := io.ReadAll(response.Body)
// 	assert.Nil(t, err)

// 	responseBody := new(model.ErrorResponse)
// 	err = json.Unmarshal(bytes, responseBody)
// 	assert.Nil(t, err)

// 	fmt.Println(response)

// 	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
// 	assert.Equal(t, http.StatusBadRequest, responseBody.Code)
// 	assert.NotEmpty(t, responseBody.Message)
// }

// func TestCreateSalesRouteIdsNotFound(t *testing.T) {
// 	defer ClearAll()

// 	token, err := GenerateTokenHelper()
// 	assert.Nil(t, err)

// 	requestBody := model.CreateSalesRequest{
// 		Name:     "Created Sales Test",
// 		Phone:    "1234560",
// 		RouteIDs: []int{1, 2},
// 	}

// 	bodyJson, err := json.Marshal(requestBody)
// 	assert.Nil(t, err)

// 	request := httptest.NewRequest(http.MethodPost, "/api/sales", strings.NewReader(string(bodyJson)))
// 	request.Header.Set("Content-Type", "application/json")
// 	request.Header.Set("Authorization", token)
// 	request.Header.Set("Accept", "application/json")

// 	response, err := app.Test(request)
// 	assert.Nil(t, err)

// 	bytes, err := io.ReadAll(response.Body)
// 	assert.Nil(t, err)

// 	responseBody := new(model.ErrorResponse)
// 	err = json.Unmarshal(bytes, responseBody)
// 	assert.Nil(t, err)

// 	fmt.Println(response)

// 	assert.Equal(t, http.StatusNotFound, response.StatusCode)
// 	assert.Equal(t, http.StatusNotFound, responseBody.Code)
// 	assert.NotEmpty(t, responseBody.Message)
// }

func TestGetAllSales(t *testing.T) {
	defer ClearAll()

	token, err := GenerateTokenHelper()
	assert.Nil(t, err)

	CreateSales(15)

	request := httptest.NewRequest(http.MethodGet, "/api/sales?page=1&perPage=10", nil)
	request.Header.Set("Authorization", token)
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[[]model.SalesResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.NotNil(t, responseBody.Data)
	assert.LessOrEqual(t, len(responseBody.Data), 10)
	assert.NotNil(t, responseBody.Paging)
	assert.Equal(t, 1, responseBody.Paging.Page)
	assert.Equal(t, 10, responseBody.Paging.PerPage)
}

func TestGetAllSalesSearch(t *testing.T) {
	defer ClearAll()

	token, err := GenerateTokenHelper()
	assert.Nil(t, err)

	CreateSales(15)

	// TODO: add search by Employee name
	request := httptest.NewRequest(http.MethodGet, "/api/sales?page=1&perPage=10&search=8", nil)
	request.Header.Set("Authorization", token)
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[[]model.SalesResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.NotNil(t, responseBody.Data)
	assert.Equal(t, len(responseBody.Data), 1)
	assert.NotNil(t, responseBody.Paging)
	assert.Equal(t, 1, responseBody.Paging.Page)
	assert.Equal(t, 10, responseBody.Paging.PerPage)
}
