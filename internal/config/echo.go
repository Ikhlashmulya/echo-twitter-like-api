package config

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

func NewEcho(config *viper.Viper) *echo.Echo {
	e := echo.New()
	e.Validator = newCustomValidator()
	e.Binder = new(validationBinder)
	e.HTTPErrorHandler = newEchoErrorHandler()

	e.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(config.GetString("jwt.secret")),
		Skipper: func(ctx echo.Context) bool {
			if (ctx.Path() == "/api/users" || ctx.Path() == "/api/users/_login") && ctx.Request().Method == http.MethodPost {
				return true
			}

			if ctx.Path() == "/api/users/:userId" && ctx.Request().Method == http.MethodGet {
				return true
			}

			if ctx.Path() == "/api/users/:userId/follow" && ctx.Request().Method == http.MethodGet {
				return true
			}

			if ctx.Path() == "/api/posts/:postId" && ctx.Request().Method == http.MethodGet {
				return true
			}

			if ctx.Path() == "/api/users/:userId/posts" && ctx.Request().Method == http.MethodGet {
				return true
			}

			if ctx.Path() == "/api/posts/:postId/comments" && ctx.Request().Method == http.MethodGet {
				return true
			}

			return false
		},
	}))

	return e
}

type validationBinder struct{}

func (vb *validationBinder) Bind(data any, ctx echo.Context) error {
	defaultBinder := new(echo.DefaultBinder)
	if err := defaultBinder.Bind(data, ctx); err != nil {
		return err
	}

	if err := ctx.Validate(data); err != nil {
		return err
	}

	return nil
}

type customValidator struct {
	validate *validator.Validate
}

func newCustomValidator() *customValidator {
	return &customValidator{validate: validator.New()}
}

func (cv *customValidator) Validate(data any) error {
	if err := cv.validate.Struct(data); err != nil {
		return err
	}
	return nil
}

func newEchoErrorHandler() echo.HTTPErrorHandler {
	return func(err error, ctx echo.Context) {
		var code int
		var message any

		if err, ok := err.(*echo.HTTPError); ok {
			code = err.Code
			message = err.Message
		}

		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errorMessages := make(map[string]any)

			for _, value := range validationErrors {
				err := value.ActualTag()
				if value.Param() != "" {
					err = fmt.Sprintf("%s=%s", err, value.Param())
				}
				errorMessages[value.Field()] = err
			}

			code = http.StatusBadRequest
			message = errorMessages
		}

		ctx.JSON(code, echo.Map{
			"errors": message,
		})
	}
}
