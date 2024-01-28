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

func (r *UserRepository) AddFollower(db *gorm.DB, userId string, user *entity.User) error {
	return db.Model(&entity.User{ID: userId}).Association("Followers").Append(user)
}

func (r *UserRepository) DeleteFollower(db *gorm.DB, userId string, user *entity.User) error {
	return db.Model(&entity.User{ID: userId}).Association("Followers").Delete(user)
}

func (r *UserRepository) FindAllFollower(db *gorm.DB, request *model.FindAllFollowerRequest) ([]entity.User, int64, error) {
	var followers []entity.User
	if err := db.Model(&entity.User{ID: request.UserID}).Limit(request.Size).Offset((request.Page - 1) * request.Size).Association("Followers").Find(&followers); err != nil {
		return nil, 0, err
	}

	total := db.Model(&entity.User{ID: request.UserID}).Association("Followers").Count()
	return followers, total, nil
}
