package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/Ikhlashmulya/echo-twitter-like-api/internal/entity"
	"github.com/Ikhlashmulya/echo-twitter-like-api/internal/model"
	"github.com/Ikhlashmulya/echo-twitter-like-api/internal/model/mapper"
	"github.com/Ikhlashmulya/echo-twitter-like-api/internal/repository"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type PostUsecase struct {
	db             *gorm.DB
	log            *logrus.Logger
	postRepository *repository.PostRepository
}

func NewPostUsecase(db *gorm.DB, log *logrus.Logger, postRepository *repository.PostRepository) *PostUsecase {
	return &PostUsecase{
		db:             db,
		log:            log,
		postRepository: postRepository,
	}
}

func (uc *PostUsecase) Create(ctx context.Context, request *model.PostCreateRequest) (*model.PostResponse, error) {
	tx := uc.db.Begin()
	defer tx.Rollback()

	post := &entity.Post{
		UserID:  request.UserID,
		Content: request.Content,
	}

	if err := uc.postRepository.Create(tx, post); err != nil {
		uc.log.Warnf("error create post: %v", err)
		return nil, echo.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		uc.log.Warnf("error commit database: %v", err)
		return nil, echo.ErrInternalServerError
	}

	return mapper.ToPostResponse(post), nil
}

func (uc *PostUsecase) Update(ctx context.Context, request *model.PostUpdateRequest) (*model.PostResponse, error) {
	tx := uc.db.Begin()
	defer tx.Rollback()

	post := new(entity.Post)
	if err := uc.postRepository.FindById(tx, post, request.ID); err != nil {
		uc.log.Warnf("error find post by id in database: %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, echo.NewHTTPError(echo.ErrNotFound.Code, fmt.Sprintf("post with id: %s not found", request.ID))
		}
		return nil, echo.ErrInternalServerError
	}

	if post.UserID != request.UserID {
		uc.log.Warnf("unauthorized user id")
		return nil, echo.ErrUnauthorized
	}

	post.Content = request.Content

	if err := uc.postRepository.Update(tx, post); err != nil {
		uc.log.Warnf("error update post: %v", err)
		return nil, echo.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		uc.log.Warnf("error commit database: %v", err)
		return nil, echo.ErrInternalServerError
	}

	return mapper.ToPostResponse(post), nil
}

func (uc *PostUsecase) Delete(ctx context.Context, request *model.PostDeleteRequest) error {
	tx := uc.db.Begin()
	defer tx.Rollback()

	post := new(entity.Post)
	if err := uc.postRepository.FindById(tx, post, request.ID); err != nil {
		uc.log.Warnf("error find post by id in database: %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(echo.ErrNotFound.Code, fmt.Sprintf("post with id: %s not found", request.ID))
		}
		return echo.ErrInternalServerError
	}

	if post.UserID != request.UserID {
		uc.log.Warnf("unauthorized user id")
		return echo.ErrUnauthorized
	}

	if err := uc.postRepository.Delete(tx, post); err != nil {
		uc.log.Warnf("error delete post: %v", err)
		return echo.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		uc.log.Warnf("error commit database: %v", err)
		return echo.ErrInternalServerError
	}

	return nil
}

func (uc *PostUsecase) FindById(ctx context.Context, postId string) (*model.PostResponse, error) {
	tx := uc.db.Begin()
	defer tx.Rollback()

	post := new(entity.Post)
	if err := uc.postRepository.FindById(tx, post, postId); err != nil {
		uc.log.Warnf("error find post by id in database: %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, echo.NewHTTPError(echo.ErrNotFound.Code, fmt.Sprintf("post with id: %s not found", postId))
		}
		return nil, echo.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		uc.log.Warnf("error commit database: %v", err)
		return nil, echo.ErrInternalServerError
	}

	return mapper.ToPostResponse(post), nil
}

func (uc *PostUsecase) FindByUserId(ctx context.Context, request *model.PostFindByUserIdRequest) (responses []model.PostResponse, total int64, err error) {
	tx := uc.db.Begin()
	defer tx.Rollback()

	posts, total, err := uc.postRepository.FindByUserId(tx, request)
	if err != nil {
		uc.log.Warnf("error find user by id from database: %v", err)
		return nil, 0, echo.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		uc.log.Warnf("error commit database: %v", err)
		return nil, 0, echo.ErrInternalServerError
	}

	for _, post := range posts {
		responses = append(responses, *mapper.ToPostResponse(&post))
	}

	return responses, total, nil
}

func (uc *PostUsecase) FindByFollowingUser(ctx context.Context, request *model.PostFindByFollowingUserRequest) (responses []model.PostResponse, total int64, err error) {
	tx := uc.db.Begin()
	defer tx.Rollback()

	posts, total, err := uc.postRepository.FindByFollowingUser(tx, request)
	if err != nil {
		uc.log.Warnf("error find by following user: %v",err)
		return nil, 0, echo.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		uc.log.Warnf("error commit database: %v", err)
		return nil, 0,echo.ErrInternalServerError
	}

	for _, post := range posts {
		responses = append(responses, *mapper.ToPostResponse(&post))
	}

	return
}
