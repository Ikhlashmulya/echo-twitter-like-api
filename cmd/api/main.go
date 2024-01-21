package main

import (
	"fmt"

	"github.com/Ikhlashmulya/echo-twitter-like-api/internal/config"
	handler "github.com/Ikhlashmulya/echo-twitter-like-api/internal/delivery/http"
)

func main() {
	configuration := config.NewViper()
	validator := config.NewValidator()
	echo := config.NewEcho(configuration, validator)
	userHandler := handler.NewUserHandler()

	echo.POST("/register", userHandler.Register)

	port := configuration.GetInt("app.port")

	echo.Logger.Fatal(echo.Start(fmt.Sprintf(":%d", port)))
}
