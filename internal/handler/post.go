package handler

import (
	"net/http"
	"strconv"

	"github.com/Ikhlashmulya/echo-twitter-like-api/internal/model"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (h *Handler) CreatePost(ctx echo.Context) (err error) {
	user := new(model.User)
	user.ID = getUserId(ctx)

	post := new(model.Post)
	post.ID = uuid.NewString()
	post.From = user.ID
	if err := ctx.Bind(post); err != nil {
		return err
	}

	if post.To == "" && post.Message == "" {
		return &echo.HTTPError{Code: echo.ErrBadGateway.Code, Message: "invalid to or message fields"}
	}

	tx := h.Db.Begin()
	defer tx.Rollback()

	if err := tx.Create(post).Error; err != nil {
		return err
	}

	tx.Commit()

	return ctx.JSON(http.StatusCreated, post)
}

func (h *Handler) FetchPost(ctx echo.Context) (err error) {
	userId := getUserId(ctx)
	page, _ := strconv.Atoi(ctx.Param("page"))
	limit, _ := strconv.Atoi(ctx.Param("limit"))

	if page == 0 {
		page = 1
	}

	if limit == 0 {
		limit = 5
	}

	tx := h.Db.Begin()
	defer tx.Rollback()

	post := []model.Post{}
	if err := tx.Limit(limit).Offset((page - 1) * limit).Find(&post, model.Post{To: userId}).Error; err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, post)
}
