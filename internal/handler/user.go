package handler

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Ikhlashmulya/echo-twitter-like-api/internal/model"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func (h *Handler) SignUp(ctx echo.Context) (err error) {
	user := &model.User{ID: uuid.NewString()}
	if err := ctx.Bind(user); err != nil {
		return err
	}

	if user.Email == "" && user.Password == "" {
		return &echo.HTTPError{Code: echo.ErrBadRequest.Code, Message: "Invalid email or password"}
	}

	tx := h.Db.Begin()
	defer tx.Rollback()

	if err := tx.Create(user).Error; err != nil {
		return err
	}

	tx.Commit()

	return ctx.JSON(http.StatusCreated, user)
}

func (h *Handler) SignIn(ctx echo.Context) (err error) {
	user := new(model.User)
	if err := ctx.Bind(user); err != nil {
		return err
	}

	tx := h.Db.Begin()
	defer tx.Rollback()

	fmt.Println(tx)

	if err := tx.Model(&model.User{}).Where("email = ? AND password = ?", user.Email, user.Password).Take(user).Error; err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return &echo.HTTPError{Code: echo.ErrBadRequest.Code, Message: "invalid email or password"}
		}
		return err
	}

	tx.Commit()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(24 * time.Hour).Unix(),
		"id": user.ID,
	})

	user.Token, err = token.SignedString([]byte("secret key"))
	if err != nil {
		return err
	}

	user.Password = ""
	return ctx.JSON(http.StatusOK, user)
}

func (h *Handler) Follow(ctx echo.Context) (err error) {
	userId := getUserId(ctx)
	id := ctx.Param("id")

	tx := h.Db.Begin()
	defer tx.Rollback()

	userFollow := new(model.User)
	if err := tx.Model(userFollow).Where("id = ?", id).Take(userFollow).Error; err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return &echo.HTTPError{Code: echo.ErrNotFound.Code, Message: "record not found"}
		}
		return err
	}

	if err := tx.Model(&model.User{ID: userId}).Association("Followers").Append(userFollow); err != nil {
		return err
	}

	tx.Commit()

	return
}

func getUserId(ctx echo.Context) string {
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	return claims["id"].(string)
}
