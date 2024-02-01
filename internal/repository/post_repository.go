package repository

import "github.com/Ikhlashmulya/echo-twitter-like-api/internal/entity"

type PostRepository struct {
	Repository[entity.Post]
}

func NewPostRepository() *PostRepository {
	return &PostRepository{}
}