package main

import (
	"github.com/Ikhlashmulya/echo-twitter-like-api/internal/handler"
	"github.com/Ikhlashmulya/echo-twitter-like-api/internal/model"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	db := newGorm()
	handler := &handler.Handler{Db: db}

	e := echo.New()
	e.Logger.SetLevel(log.ERROR)
	e.Use(middleware.Logger())
	e.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte("secret key"),
		Skipper: func(c echo.Context) bool {
			if c.Path() == "/login" || c.Path() == "/signup" {
				return true
			}
			return false
		},
	}))

	e.POST("/login", handler.SignIn)
	e.POST("/signup", handler.SignUp)
	e.POST("/follow/:id", handler.Follow)
	e.POST("/posts", handler.CreatePost)
	e.GET("/feed", handler.FetchPost)

	e.Logger.Fatal(e.Start(":9000"))
}

func newGorm() *gorm.DB {
	dsn := "root:@tcp(127.0.0.1:3306)/twitter?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.Post{})

	return db
}
