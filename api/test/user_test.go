package test

import (
	// "context"
	"api/internal/model"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	defer ClearAll()
	requestBody := model.RegisterUserRequest{
		Name:     "Rafli",
		Username: "rafli@gmail.com",
		Password: "rahasia",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPost, "/api/users", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, requestBody.Name, responseBody.Data.Name)
	assert.Equal(t, requestBody.Username, responseBody.Data.Username)
	assert.NotNil(t, responseBody.Data.ID)
}

func TestFindAll(t *testing.T) {
	defer ClearAll()
	CreateUsers(20)
	request := httptest.NewRequest(http.MethodGet, "/api/users", nil)
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[[]model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, 20, len(responseBody.Data))
	assert.Nil(t, responseBody.Paging)
}

func TestFindAllWithPagination(t *testing.T) {
	defer ClearAll()
	CreateUsers(20)

	request := httptest.NewRequest(http.MethodGet, "/api/users?page=1&perPage=10", nil)
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[[]model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, 10, len(responseBody.Data))
	assert.Equal(t, int64(20), responseBody.Paging.TotalItem)
	assert.Equal(t, int64(2), responseBody.Paging.TotalPage)
	assert.Equal(t, 1, responseBody.Paging.Page)
	assert.Equal(t, 10, responseBody.Paging.PerPage)
}

func TestUpdate(t *testing.T) {
	defer ClearAll()
	user := CreateUsers(1)[0]

	requestBody := model.UpdateUserRequest{
		Name:     "Rafli edit",
		Username: "rafli999@gmail.com",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPut, "/api/users/"+user.ID.String(), strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, user.ID, responseBody.Data.ID)
	assert.Equal(t, requestBody.Name, responseBody.Data.Name)
	assert.Equal(t, requestBody.Username, responseBody.Data.Username)
}

func TestDelete(t *testing.T) {
	defer ClearAll()
	user := CreateUsers(1)[0]

	request := httptest.NewRequest(http.MethodDelete, "/api/users/"+user.ID.String(), nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[bool])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, responseBody.Data, true)
}
