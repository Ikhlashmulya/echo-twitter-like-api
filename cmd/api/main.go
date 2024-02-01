package main

import (
	"fmt"

	"github.com/Ikhlashmulya/echo-twitter-like-api/internal/config"
	handler "github.com/Ikhlashmulya/echo-twitter-like-api/internal/delivery/http"
	"github.com/Ikhlashmulya/echo-twitter-like-api/internal/repository"
	"github.com/Ikhlashmulya/echo-twitter-like-api/internal/usecase"
)

func main() {
	configuration := config.NewViper()
	log := config.NewLogger(configuration)
	db := config.NewGorm(configuration, log)

	userRepository := repository.NewUserRepository()
	userUsecase := usecase.NewUserUsecase(db, log, configuration, userRepository)
	userHandler := handler.NewUserHandler(log, userUsecase)

	postRepository := repository.NewPostRepository()
	postUsecase := usecase.NewPostUsecase(db, log, postRepository)
	postHandler := handler.NewPostHandler(log, postUsecase)
	
	echo := config.NewEcho(configuration)

	api := echo.Group("/api")
	userHandler.Route(api)
	postHandler.Route(api)

	port := configuration.GetInt("app.port")

	echo.Logger.Fatal(echo.Start(fmt.Sprintf(":%d", port)))
}
