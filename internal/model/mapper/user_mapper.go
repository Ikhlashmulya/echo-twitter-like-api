package mapper

import (
	"github.com/Ikhlashmulya/echo-twitter-like-api/internal/entity"
	"github.com/Ikhlashmulya/echo-twitter-like-api/internal/model"
)

func ToUserResponse(user *entity.User) *model.UserResponse {
	return &model.UserResponse{
		ID:           user.ID,
		Name:         user.Name,
		Email:        user.Email,
		PhotoProfile: user.PhotoProfile,
	}
}

func ToUserTokenResponse(token string) *model.UserTokenResponse {
	return &model.UserTokenResponse{
		Token: token,
	}
}
