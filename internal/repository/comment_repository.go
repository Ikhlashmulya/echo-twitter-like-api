package repository

import (
	"github.com/Ikhlashmulya/echo-twitter-like-api/internal/entity"
	"github.com/Ikhlashmulya/echo-twitter-like-api/internal/model"
	"gorm.io/gorm"
)

type CommentRepository struct {
	Repository[entity.Comment]
}

func NewCommentRepository() *CommentRepository {
	return &CommentRepository{}
}

func (r *CommentRepository) FindByPostId(tx *gorm.DB, request *model.CommentFindByPostId) ([]entity.Comment, int64, error) {
	var comments []entity.Comment
	if err := tx.Find(&comments, "post_id = ?", request.PostID).Limit(request.Size).Offset((request.Page-1)*request.Size).Error; err != nil {
		return nil, 0, err
	}

	total := tx.Model(&entity.Comment{PostID: request.PostID}).Association("Post").Count()

	return comments, total, nil
}
