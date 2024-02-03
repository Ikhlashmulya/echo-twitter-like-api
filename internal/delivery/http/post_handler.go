package http

import (
	"net/http"

	"github.com/Ikhlashmulya/echo-twitter-like-api/internal/delivery/http/util"
	"github.com/Ikhlashmulya/echo-twitter-like-api/internal/model"
	"github.com/Ikhlashmulya/echo-twitter-like-api/internal/usecase"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type PostHandler struct {
	log         *logrus.Logger
	postUsecase *usecase.PostUsecase
}

func NewPostHandler(log *logrus.Logger, postUsecase *usecase.PostUsecase) *PostHandler {
	return &PostHandler{
		log:         log,
		postUsecase: postUsecase,
	}
}

func (h *PostHandler) Route(router *echo.Group) {
	router.POST("/posts", h.Create)
	router.GET("/posts/:postId", h.FindById)
	router.PUT("/posts/:postId", h.Update)
	router.DELETE("/posts/:postId", h.Delete)
	router.GET("/users/:userId/posts", h.FindByUserId)
}

func (h *PostHandler) Create(ctx echo.Context) error {
	request := new(model.PostCreateRequest)
	if err := ctx.Bind(request); err != nil {
		h.log.Warnf("error binding from request body: %v", err)
		return err
	}

	request.UserID = util.GetAuthId(ctx)

	response, err := h.postUsecase.Create(ctx.Request().Context(), request)
	if err != nil {
		h.log.Warnf("error create post: %v", err)
		return err
	}

	return ctx.JSON(http.StatusCreated, &model.WebResponse[*model.PostResponse]{Data: response})
}

func (h *PostHandler) Update(ctx echo.Context) error {
	request := new(model.PostUpdateRequest)
	if err := ctx.Bind(request); err != nil {
		h.log.Warnf("error binding from request body: %v", err)
		return err
	}

	request.UserID = util.GetAuthId(ctx)

	response, err := h.postUsecase.Update(ctx.Request().Context(), request)
	if err != nil {
		h.log.Warnf("error create post: %v", err)
		return err
	}

	return ctx.JSON(http.StatusOK, &model.WebResponse[*model.PostResponse]{Data: response})
}

func (h *PostHandler) Delete(ctx echo.Context) error {
	request := new(model.PostDeleteRequest)
	request.ID = ctx.Param("postId")
	request.UserID = util.GetAuthId(ctx)

	if err := h.postUsecase.Delete(ctx.Request().Context(), request); err != nil {
		h.log.Warnf("error delete post: %v", err)
		return err
	}

	return ctx.JSON(http.StatusOK, &model.WebResponse[string]{Data: "OK"})
}

func (h *PostHandler) FindById(ctx echo.Context) error {
	postId := ctx.Param("postId")

	response, err := h.postUsecase.FindById(ctx.Request().Context(), postId)
	if err != nil {
		h.log.Warnf("error find post by id: %v", err)
		return err
	}

	return ctx.JSON(http.StatusOK, &model.WebResponse[*model.PostResponse]{Data: response})
}

func (h *PostHandler) FindByUserId(ctx echo.Context) error {
	request := new(model.PostFindByUserIdRequest)
	if err := ctx.Bind(request); err != nil {
		h.log.Warnf("error binding : %v", err)
		return err
	}

	request.UserID = ctx.Param("userId")

	if request.Page == 0 {
		request.Page = 1
	}

	if request.Size == 0 {
		request.Size = 5
	}

	responses, total, err := h.postUsecase.FindByUserId(ctx.Request().Context(), request)
	if err != nil {
		h.log.Warnf("error find post by user id: %v", err)
		return err
	}

	return ctx.JSON(http.StatusOK, &model.WebResponse[[]model.PostResponse]{
		Data: responses,
		Paging: &model.Paging{
			Page:  request.Page,
			Size:  request.Size,
			Total: total,
		},
	})
}
