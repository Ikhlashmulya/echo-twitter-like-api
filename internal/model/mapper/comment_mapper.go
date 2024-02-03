package mapper

import (
	"github.com/Ikhlashmulya/echo-twitter-like-api/internal/entity"
	"github.com/Ikhlashmulya/echo-twitter-like-api/internal/model"
)

func ToCommentResponse(comment *entity.Comment) *model.CommentResponse {
	return &model.CommentResponse{
		ID:        comment.ID,
		User:      comment.UserID,
		Reply:     comment.Reply,
		CreatedAt: comment.CreatedAt,
		UpdatedAt: comment.UpdatedAt,
	}
}
