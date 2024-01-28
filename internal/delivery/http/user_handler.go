package http

import (
	"net/http"

	"github.com/Ikhlashmulya/echo-twitter-like-api/internal/helper"
	"github.com/Ikhlashmulya/echo-twitter-like-api/internal/model"
	"github.com/Ikhlashmulya/echo-twitter-like-api/internal/usecase"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type UserHandler struct {
	log         *logrus.Logger
	userUsecase *usecase.UserUsecase
}

func NewUserHandler(log *logrus.Logger, userUsecase *usecase.UserUsecase) *UserHandler {
	return &UserHandler{
		log:         log,
		userUsecase: userUsecase,
	}
}

func (h *UserHandler) Route(route *echo.Group) {
	route.POST("/user", h.Register)
	route.POST("/user/_login", h.Login)
	route.GET("/user/:userId", h.FindById)
	route.POST("/user/:userId/follow", h.AddFollower)
	route.DELETE("/user/:userId/follow", h.DeleteFollower)
	route.GET("/user/:userId/follow", h.FindAllFollower)
}

func (h *UserHandler) Register(ctx echo.Context) (err error) {
	request := new(model.UserRegisterRequest)
	if err := ctx.Bind(request); err != nil {
		h.log.Warnf("error binding from request body: %v", err)
		return err
	}

	response, err := h.userUsecase.Register(ctx.Request().Context(), request)
	if err != nil {
		h.log.Warnf("error register user: %v", err)
		return err
	}

	return ctx.JSON(http.StatusOK, &model.WebResponse[*model.UserResponse]{Data: response})
}

func (h *UserHandler) Login(ctx echo.Context) (err error) {
	request := new(model.UserLoginRequest)
	if err := ctx.Bind(request); err != nil {
		h.log.Warnf("error binding from request body: %v", err)
		return err
	}

	response, err := h.userUsecase.Login(ctx.Request().Context(), request)
	if err != nil {
		h.log.Warnf("error login user: %v", err)
		return err
	}

	return ctx.JSON(http.StatusOK, &model.WebResponse[*model.UserTokenResponse]{Data: response})
}

func (h *UserHandler) FindById(ctx echo.Context) error {
	userId := ctx.Param("userId")

	response, err := h.userUsecase.FindById(ctx.Request().Context(), userId)
	if err != nil {
		h.log.Warnf("error find user by id: %v", err)
		return err
	}

	return ctx.JSON(http.StatusOK, &model.WebResponse[*model.UserResponse]{Data: response})
}

func (h *UserHandler) AddFollower(ctx echo.Context) error {
	authId := helper.GetAuthId(ctx)

	userId := ctx.Param("userId")

	if err := h.userUsecase.AddFollower(ctx.Request().Context(), userId, authId); err != nil {
		h.log.Warnf("error add follower: %v", err)
		return err
	}

	return ctx.JSON(http.StatusOK, &model.WebResponse[string]{Data: "OK"})
}

func (h *UserHandler) DeleteFollower(ctx echo.Context) error {
	authId := helper.GetAuthId(ctx)

	userId := ctx.Param("userId")

	if err := h.userUsecase.DeleteFollower(ctx.Request().Context(), userId, authId); err != nil {
		h.log.Warnf("error delete follower: %v", err)
		return err
	}

	return ctx.JSON(http.StatusOK, &model.WebResponse[string]{Data: "OK"})
}

func (h *UserHandler) FindAllFollower(ctx echo.Context) error {
	request := new(model.FindAllFollowerRequest)
	if err := ctx.Bind(request); err != nil {
		h.log.Warnf("error binding request: %v", err)
		return err
	}

	request.UserID = ctx.Param("userId")

	if request.Page == 0 {
		request.Page = 1
	}

	if request.Size == 0 {
		request.Size = 5
	}

	responses, total, err := h.userUsecase.FindAllFollower(ctx.Request().Context(), request)
	if err != nil {
		h.log.Warnf("error find all follower: %v", err)
		return err
	}

	return ctx.JSON(http.StatusOK, &model.WebResponse[[]model.UserResponse]{
		Data: responses,
		Paging: &model.Paging{
			Page:  request.Page,
			Size:  request.Size,
			Total: total,
		},
	})
}
