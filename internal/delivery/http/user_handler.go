package http

import (
	"net/http"

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

func (h *UserHandler) Register(ctx echo.Context) (err error) {
	request := new(model.UserRegisterRequest)
	if err := ctx.Bind(request); err != nil {
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
		return err
	}

	response, err := h.userUsecase.Login(ctx.Request().Context(), request)
	if err != nil {
		h.log.Warnf("error login user: %v", err)
		return err
	}

	return ctx.JSON(http.StatusOK, &model.WebResponse[*model.UserTokenResponse]{Data: response})
}