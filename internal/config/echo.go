package config

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

func NewEcho(config *viper.Viper) *echo.Echo {
	echo := echo.New()
	echo.Validator = newCustomValidator()
	echo.Binder = new(ValidationBinder)
	echo.HTTPErrorHandler = newEchoErrorHandler()

	return echo
}

type ValidationBinder struct{}

func (vb *ValidationBinder) Bind(data any, ctx echo.Context) error {
	defaultBinder := new(echo.DefaultBinder)
	if err := defaultBinder.Bind(data, ctx); err != nil {
		return err
	}

	if err := ctx.Validate(data); err != nil {
		return err
	}

	return nil
}

type CustomValidator struct {
	validate *validator.Validate
}

func newCustomValidator() *CustomValidator {
	return &CustomValidator{validate: validator.New()}
}

func (cv *CustomValidator) Validate(data any) error {
	if err := cv.validate.Struct(data); err != nil {
		return err
	}
	return nil
}

func newEchoErrorHandler() echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
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

		c.JSON(code, echo.Map{
			"errors": message,
		})
	}
}
