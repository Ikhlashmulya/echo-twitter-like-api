package repository

import (
	"github.com/Ikhlashmulya/echo-twitter-like-api/internal/entity"
	"github.com/Ikhlashmulya/echo-twitter-like-api/internal/model"
	"gorm.io/gorm"
)

type UserRepository struct {
	Repository[entity.User]
}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) AddFollowing(db *gorm.DB, authId string, user *entity.User) error {
	return db.Model(&entity.User{ID: authId}).Association("Following").Append(user)
}

func (r *UserRepository) DeleteFollowing(db *gorm.DB, authId string, user *entity.User) error {
	return db.Model(&entity.User{ID: authId}).Association("Following").Delete(user)
}

func (r *UserRepository) FindAllFollowing(db *gorm.DB, request *model.UserFindAllFollowingRequest) ([]entity.User, int64, error) {
	var followers []entity.User
	if err := db.Model(&entity.User{ID: request.UserID}).Limit(request.Size).Offset((request.Page - 1) * request.Size).Association("Following").Find(&followers); err != nil {
		return nil, 0, err
	}

	total := db.Model(&entity.User{ID: request.UserID}).Association("Following").Count()
	return followers, total, nil
}
