package test

import (
	"api/internal/entity/enum"
	"api/internal/model"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// ==================== AUTH TESTS ====================

func TestLogin(t *testing.T) {
	defer ClearAll()

	user := CreateUsers(1)[0]

	requestBody := model.LoginUserRequest{
		Username: user.Username,
		Password: "rahasia",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPost, "/api/login", strings.NewReader(string(bodyJson)))
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
	assert.NotEmpty(t, responseBody.Token)
	assert.Equal(t, user.ID, responseBody.Data.ID)
	assert.Equal(t, user.Username, responseBody.Data.Username)
	assert.Equal(t, user.Name, responseBody.Data.Name)
}

func TestLoginInvalidCredentials(t *testing.T) {
	defer ClearAll()

	CreateUsers(1)

	requestBody := model.LoginUserRequest{
		Username: "wronguser",
		Password: "wrongpass",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPost, "/api/login", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
}

func TestGetCurrentUser(t *testing.T) {
	defer ClearAll()

	token, err := GenerateTokenHelper()
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodGet, "/api/current", nil)
	request.Header.Set("Authorization", token)
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.NotNil(t, responseBody.Data.ID)
	assert.NotEmpty(t, responseBody.Data.Username)
	assert.NotEmpty(t, responseBody.Data.Name)
}

func TestGetCurrentUserUnauthorized(t *testing.T) {
	defer ClearAll()

	request := httptest.NewRequest(http.MethodGet, "/api/current", nil)
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
}

// ==================== USER TESTS ====================

func TestCreateUser(t *testing.T) {
	defer ClearAll()

	token, err := GenerateTokenHelper()
	assert.Nil(t, err)

	requestBody := model.RegisterUserRequest{
		Name:     "Created User Test",
		Username: "username123",
		Password: "rahasia",
		Phone:    "08123456789",
		Role:     enum.SUPER_ADMIN,
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPost, "/api/users", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", token)
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusCreated, response.StatusCode)
	assert.Equal(t, requestBody.Name, responseBody.Data.Name)
	assert.Equal(t, requestBody.Username, responseBody.Data.Username)
	assert.Equal(t, requestBody.Phone, responseBody.Data.Phone)
	assert.Equal(t, requestBody.Role, responseBody.Data.Role)
	assert.NotNil(t, responseBody.Data.ID)
	assert.NotEmpty(t, responseBody.Data.CreatedAt)
}

func TestCreateUserUnauthorized(t *testing.T) {
	defer ClearAll()

	requestBody := model.RegisterUserRequest{
		Name:     "Test User",
		Username: "testuser",
		Password: "rahasia",
		Role:     enum.SUPER_ADMIN,
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPost, "/api/users", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
}

func TestCreateUserInvalidRequest(t *testing.T) {
	defer ClearAll()

	token, err := GenerateTokenHelper()
	assert.Nil(t, err)

	requestBody := model.RegisterUserRequest{
		Name:     "",
		Username: "",
		Password: "",
		Phone:    "",
		Role:     enum.OWNER,
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPost, "/api/users", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", token)
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
}

func TestGetAllUsers(t *testing.T) {
	defer ClearAll()

	token, err := GenerateTokenHelper()
	assert.Nil(t, err)

	CreateUsers(15)

	request := httptest.NewRequest(http.MethodGet, "/api/users?page=1&perPage=10", nil)
	request.Header.Set("Authorization", token)
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[[]model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.NotNil(t, responseBody.Data)
	assert.LessOrEqual(t, len(responseBody.Data), 10)
	assert.NotNil(t, responseBody.Paging)
	assert.Equal(t, 1, responseBody.Paging.Page)
	assert.Equal(t, 10, responseBody.Paging.PerPage)
}

func TestGetAllUsersWithFilters(t *testing.T) {
	defer ClearAll()

	token, err := GenerateTokenHelper()
	assert.Nil(t, err)

	users := CreateUsers(5)

	request := httptest.NewRequest(http.MethodGet, "/api/users?search="+users[0].Name, nil)
	request.Header.Set("Authorization", token)
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[[]model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, len(responseBody.Data), 0)
}

func TestGetAllUsersUnauthorized(t *testing.T) {
	defer ClearAll()

	request := httptest.NewRequest(http.MethodGet, "/api/users", nil)
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
}

func TestUpdateUser(t *testing.T) {
	defer ClearAll()

	token, err := GenerateTokenHelper()
	assert.Nil(t, err)

	user := CreateUsers(1)[0]

	requestBody := model.RegisterUserRequest{
		Name:     "Updated Name",
		Username: "updatedusername",
		Password: "newpassword",
		Phone:    "08199999999",
		Role:     enum.OWNER,
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPut, "/api/users/"+user.ID.String(), strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", token)
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
	assert.Equal(t, requestBody.Phone, responseBody.Data.Phone)
	assert.Equal(t, requestBody.Role, responseBody.Data.Role)
}

func TestUpdateUserNotFound(t *testing.T) {
	defer ClearAll()

	token, err := GenerateTokenHelper()
	assert.Nil(t, err)

	requestBody := model.RegisterUserRequest{
		Name:     "Updated Name",
		Username: "updated",
		Password: "password",
		Role:     enum.OWNER,
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	fakeID := uuid.New().String()
	request := httptest.NewRequest(http.MethodPut, "/api/users/"+fakeID, strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", token)
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusNotFound, response.StatusCode)
}

func TestDeleteUser(t *testing.T) {
	defer ClearAll()

	token, err := GenerateTokenHelper()
	assert.Nil(t, err)

	user := CreateUsers(1)[0]

	request := httptest.NewRequest(http.MethodDelete, "/api/users/"+user.ID.String(), nil)
	request.Header.Set("Authorization", token)
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[bool])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.True(t, responseBody.Data)
}

func TestDeleteUserUnauthorized(t *testing.T) {
	defer ClearAll()

	user := CreateUsers(1)[0]

	request := httptest.NewRequest(http.MethodDelete, "/api/users/"+user.ID.String(), nil)
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
}
