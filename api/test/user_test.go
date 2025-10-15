package test

import (
	// "context"
	"api/internal/model"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	// "time"

	// "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	defer ClearAll()
	requestBody := model.RegisterUserRequest{
		Name:     "Rafli",
		Email:    "rafli@gmail.com",
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
	assert.Equal(t, requestBody.Email, responseBody.Data.Email)
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
		Name:  "Rafli edit",
		Email: "rafli999@gmail.com",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPut, "/api/users/"+strconv.Itoa(user.ID), strings.NewReader(string(bodyJson)))
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
	assert.Equal(t, requestBody.Email, responseBody.Data.Email)
}

func TestDelete(t *testing.T) {
	defer ClearAll()
	user := CreateUsers(1)[0]

	request := httptest.NewRequest(http.MethodDelete, "/api/users/"+strconv.Itoa(user.ID), nil)
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

// func TestRedis(t *testing.T) {
// 	rdb := redis.NewClient(&redis.Options{
// 		Addr:             "localhost:6379",
// 		Password:         "",
// 		DB:               0,
// 		DisableIndentity: true,
// 	})
// 	ctx := context.Background()

// 	t.Run("conection", func(t *testing.T) {
// 		assert.NotNil(t, rdb)

// 		// err := rdb.Close()
// 		// assert.Nil(t, err)
// 	})

// 	t.Run("Ping", func(t *testing.T) {
// 		result, err := rdb.Ping(ctx).Result()

// 		assert.Nil(t, err)
// 		assert.Equal(t, "PONG", result)
// 	})

// 	t.Run("Set", func(t *testing.T) {
// 		rdb.SetEx(ctx, "key", "value", time.Second)

// 		result, err := rdb.Get(ctx, "key").Result()
// 		assert.Nil(t, err)

// 		assert.Equal(t, "value", result)
// 	})

// 	t.Run("List", func(t *testing.T) {
// 		rdb.RPush(ctx, "key", "value")
// 		rdb.LPush(ctx, "key", "value left")
// 		rdb.RPush(ctx, "key", "value right")

// 		assert.Equal(t, "value right", rdb.RPop(ctx, "key").Val())
// 		assert.Equal(t, "value left", rdb.LPop(ctx, "key").Val())
// 		assert.Equal(t, "value", rdb.LPop(ctx, "key").Val())

// 		rdb.Del(ctx, "key")
// 	})

// }
