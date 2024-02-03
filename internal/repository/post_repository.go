package repository

import (
	"github.com/Ikhlashmulya/echo-twitter-like-api/internal/entity"
	"github.com/Ikhlashmulya/echo-twitter-like-api/internal/model"
	"gorm.io/gorm"
)

type PostRepository struct {
	Repository[entity.Post]
}

func NewPostRepository() *PostRepository {
	return &PostRepository{}
}

func (r *PostRepository) FindByUserId(tx *gorm.DB, request *model.PostFindByUserIdRequest) ([]entity.Post, int64, error) {
	var posts []entity.Post
	if err := tx.Where("user_id = ?", request.UserID).Limit(request.Size).Offset((request.Page - 1) * request.Size).Find(&posts).Error; err != nil {
		return nil, 0, err
	}

	var total int64
	if err := tx.Model(new(entity.Post)).Where("user_id = ?", request.UserID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}
