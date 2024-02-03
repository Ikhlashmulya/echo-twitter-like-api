package http

import (
	"net/http"

	"github.com/Ikhlashmulya/echo-twitter-like-api/internal/delivery/http/util"
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

func (h *UserHandler) Route(router *echo.Group) {
	router.POST("/users", h.Register)
	router.POST("/users/_login", h.Login)
	router.GET("/users/:userId", h.FindById)
	router.POST("/users/:userId/follow", h.AddFollowing)
	router.DELETE("/users/:userId/follow", h.DeleteFollowing)
	router.GET("/users/:userId/follow", h.FindAllFollowing)
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

func (h *UserHandler) AddFollowing(ctx echo.Context) error {
	authId := util.GetAuthId(ctx)

	userId := ctx.Param("userId")

	if err := h.userUsecase.AddFollowing(ctx.Request().Context(), userId, authId); err != nil {
		h.log.Warnf("error add follower: %v", err)
		return err
	}

	return ctx.JSON(http.StatusOK, &model.WebResponse[string]{Data: "OK"})
}

func (h *UserHandler) DeleteFollowing(ctx echo.Context) error {
	authId := util.GetAuthId(ctx)

	userId := ctx.Param("userId")

	if err := h.userUsecase.DeleteFollowing(ctx.Request().Context(), userId, authId); err != nil {
		h.log.Warnf("error delete follower: %v", err)
		return err
	}

	return ctx.JSON(http.StatusOK, &model.WebResponse[string]{Data: "OK"})
}

func (h *UserHandler) FindAllFollowing(ctx echo.Context) error {
	request := new(model.UserFindAllFollowingRequest)
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

	responses, total, err := h.userUsecase.FindAllFollowing(ctx.Request().Context(), request)
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
