package config

import (
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

func NewEcho(config *viper.Viper, validator *Validator) *echo.Echo {
	echo := echo.New()
	echo.Validator = validator
	echo.Binder = new(ValidationBinder)

	return echo
}

type ValidationBinder struct{}

func (vb *ValidationBinder) Bind(data any, ctx echo.Context) (err error) {
	defaultBinder := new(echo.DefaultBinder)
	if err = defaultBinder.Bind(data, ctx); err != nil {
		return
	}

	if err = ctx.Validate(data); err != nil {
		return
	}

	return
}
