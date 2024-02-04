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

type CommentUsecase struct {
	db                *gorm.DB
	log               *logrus.Logger
	commentRepository *repository.CommentRepository
}

func NewCommentUsecase(db *gorm.DB, log *logrus.Logger, commentRepository *repository.CommentRepository) *CommentUsecase {
	return &CommentUsecase{
		db:                db,
		log:               log,
		commentRepository: commentRepository,
	}
}

func (uc *CommentUsecase) Create(ctx context.Context, request *model.CommentCreateRequest) (*model.CommentResponse, error) {
	tx := uc.db.Begin()
	defer tx.Rollback()

	comment := &entity.Comment{
		UserID: request.UserID,
		PostID: request.PostID,
		Reply:  request.Reply,
	}

	if err := uc.commentRepository.Create(tx, comment); err != nil {
		uc.log.Warnf("error create comment in database: %v", err)
		return nil, echo.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		uc.log.Warnf("error commit database: %v", err)
		return nil, echo.ErrInternalServerError
	}

	return mapper.ToCommentResponse(comment), nil
}

func (uc *CommentUsecase) Update(ctx context.Context, request *model.CommentUpdateRequest) (*model.CommentResponse, error) {
	tx := uc.db.Begin()
	defer tx.Rollback()

	comment := new(entity.Comment)
	if err := uc.commentRepository.FindById(tx, comment, request.ID); err != nil {
		uc.log.Warnf("error find comment by id: %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, echo.NewHTTPError(echo.ErrNotFound.Code, fmt.Sprintf("comment with id: %s not found", request.ID))
		}
		return nil, echo.ErrInternalServerError
	}

	if comment.UserID != request.UserID {
		uc.log.Warnf("unauthorized user id")
		return nil, echo.ErrUnauthorized
	}

	comment.Reply = request.Reply

	if err := uc.commentRepository.Update(tx, comment); err != nil {
		uc.log.Warnf("error update comment by id: %v", err)
		return nil, echo.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		uc.log.Warnf("error commit database: %v", err)
		return nil, echo.ErrInternalServerError
	}

	return mapper.ToCommentResponse(comment), nil
}

func (uc *CommentUsecase) Delete(ctx context.Context, request *model.CommentDeleteRequest) error {
	tx := uc.db.Begin()
	defer tx.Rollback()

	comment := new(entity.Comment)
	if err := uc.commentRepository.FindById(tx, comment, request.ID); err != nil {
		uc.log.Warnf("error find comment by id: %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(echo.ErrNotFound.Code, fmt.Sprintf("comment with id: %s not found", request.ID))
		}
		return echo.ErrInternalServerError
	}

	if comment.UserID != request.UserID {
		uc.log.Warnf("unauthorized user id")
		return echo.ErrUnauthorized
	}

	if err := uc.commentRepository.Delete(tx, comment); err != nil {
		uc.log.Warnf("error delete comment by id: %v", err)
		return echo.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		uc.log.Warnf("error commit database: %v", err)
		return echo.ErrInternalServerError
	}

	return nil
}

func (uc *CommentUsecase) FindByPostId(ctx context.Context, request *model.CommentFindByPostId) (responses []model.CommentResponse, total int64, err error) {
	tx := uc.db.Begin()
	defer tx.Rollback()

	comments, total, err := uc.commentRepository.FindByPostId(tx, request)
	if err != nil {
		uc.log.Warnf("error find comment by post id: %v", err)
		return nil, 0, echo.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		uc.log.Warnf("error commit database: %v", err)
		return nil, 0, echo.ErrInternalServerError
	}

	for _, comment := range comments {
		responses = append(responses, *mapper.ToCommentResponse(&comment))
	}

	return
}
