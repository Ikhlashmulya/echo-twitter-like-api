package repository

import (
	"github.com/Ikhlashmulya/echo-twitter-like-api/internal/entity"
)

type UserRepository struct {
	Repository[entity.User]
}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}
