package mapper

import (
	"github.com/Ikhlashmulya/echo-twitter-like-api/internal/entity"
	"github.com/Ikhlashmulya/echo-twitter-like-api/internal/model"
)

func ToPostResponse(post *entity.Post) *model.PostResponse {
	return &model.PostResponse{
		ID:        post.ID,
		UserID:    post.UserID,
		Content:   post.Content,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
	}
}
