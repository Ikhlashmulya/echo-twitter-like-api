package http

import (
	"net/http"

	"github.com/Ikhlashmulya/echo-twitter-like-api/internal/delivery/http/util"
	"github.com/Ikhlashmulya/echo-twitter-like-api/internal/model"
	"github.com/Ikhlashmulya/echo-twitter-like-api/internal/usecase"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type CommentHandler struct {
	log            *logrus.Logger
	commentUsecase *usecase.CommentUsecase
}

func NewCommentHandler(log *logrus.Logger, commentUsecase *usecase.CommentUsecase) *CommentHandler {
	return &CommentHandler{
		log:            log,
		commentUsecase: commentUsecase,
	}
}

func (h *CommentHandler) Route(router *echo.Group) {
	router.POST("/posts/:postId/comments", h.Create)
	router.GET("/posts/:postId/comments", h.FindByPostId)
	router.PUT("/comments/:commentId", h.Update)
	router.DELETE("/comments/:commentId", h.Delete)
}

func (h *CommentHandler) Create(ctx echo.Context) error {
	request := new(model.CommentCreateRequest)
	if err := ctx.Bind(request); err != nil {
		h.log.Warnf("error binding from request body: %v", err)
		return err
	}

	request.UserID = util.GetAuthId(ctx)
	request.PostID = ctx.Param("postId")

	response, err := h.commentUsecase.Create(ctx.Request().Context(), request)
	if err != nil {
		h.log.Warnf("error create comment: %v", err)
		return err
	}

	return ctx.JSON(http.StatusOK, &model.WebResponse[*model.CommentResponse]{
		Data: response,
	})
}

func (h *CommentHandler) Update(ctx echo.Context) error {
	request := new(model.CommentUpdateRequest)
	if err := ctx.Bind(request); err != nil {
		h.log.Warnf("error binding from request body: %v", err)
		return err
	}

	request.UserID = util.GetAuthId(ctx)
	request.ID = ctx.Param("commentId")

	response, err := h.commentUsecase.Update(ctx.Request().Context(), request)
	if err != nil {
		h.log.Warnf("error update comment: %v", err)
		return err
	}

	return ctx.JSON(http.StatusOK, &model.WebResponse[*model.CommentResponse]{
		Data: response,
	})
}

func (h *CommentHandler) Delete(ctx echo.Context) error {
	request := new(model.CommentDeleteRequest)
	request.UserID = util.GetAuthId(ctx)
	request.ID = ctx.Param("commentId")

	if err := h.commentUsecase.Delete(ctx.Request().Context(), request); err != nil {
		h.log.Warnf("error delete comment: %v", err)
		return err
	}

	return ctx.JSON(http.StatusOK, &model.WebResponse[string]{
		Data: "OK",
	})
}

func (h *CommentHandler) FindByPostId(ctx echo.Context) error {
	request := new(model.CommentFindByPostId)
	if err := ctx.Bind(request); err != nil {
		h.log.Warnf("error binding: %v", err)
		return err
	}

	if request.Page == 0 {
		request.Page = 1
	}

	if request.Size == 0 {
		request.Size = 5
	}

	request.PostID = ctx.Param("postId")

	responses, total, err := h.commentUsecase.FindByPostId(ctx.Request().Context(), request)
	if err != nil {
		h.log.Warnf("error find comment: %v", err)
		return err
	}

	return ctx.JSON(http.StatusOK, &model.WebResponse[[]model.CommentResponse]{
		Data: responses,
		Paging: &model.Paging{
			Total: total,
			Page:  request.Page,
			Size:  request.Size,
		},
	})
}
