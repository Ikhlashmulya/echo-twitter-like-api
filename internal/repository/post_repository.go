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

func (r *PostRepository) FindByFollowingUser(tx *gorm.DB, request *model.PostFindByFollowingUserRequest) ([]entity.Post, int64, error) {
	var results []entity.Post
	tx.Raw("SELECT id, user_id, content, created_at, updated_at FROM posts WHERE user_id IN (SELECT following_id FROM user_followers WHERE user_id = ?) ORDER BY created_at DESC LIMIT ? OFFSET ?", request.UserID, request.Size, (request.Page-1)*request.Size).Scan(&results)

	var total int64
	tx.Raw("SELECT count(id) FROM posts WHERE user_id IN (SELECT following_id FROM user_followers WHERE user_id = ?) ORDER BY max(created_at) DESC", request.UserID).Scan(&total)

	return results, total, nil
}
