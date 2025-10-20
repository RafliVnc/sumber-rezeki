package converter

import (
	"api/internal/entity"
	"api/internal/model"
)

func ToUserResponse(user *entity.User) *model.UserResponse {
	return &model.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Username:  user.Username,
		Phone:     user.Phone,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	}
}

func UserToTokenResponse(user *entity.User, token string) *model.LoginUserResponse {
	return &model.LoginUserResponse{
		User:  *ToUserResponse(user),
		Token: token,
	}
}
