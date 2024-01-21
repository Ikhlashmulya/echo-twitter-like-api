package http

import (
	"fmt"

	"github.com/Ikhlashmulya/echo-twitter-like-api/internal/model"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

func (h *UserHandler) Register(ctx echo.Context) (err error) {
	request := new(model.UserRegisterRequest)
	if err := ctx.Bind(request); err != nil {
		return err
	}

	fmt.Println(request)

	return ctx.String(200, "OK")
}
