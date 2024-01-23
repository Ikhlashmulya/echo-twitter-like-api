package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Ikhlashmulya/echo-twitter-like-api/internal/entity"
	"github.com/Ikhlashmulya/echo-twitter-like-api/internal/model"
	"github.com/Ikhlashmulya/echo-twitter-like-api/internal/model/mapper"
	"github.com/Ikhlashmulya/echo-twitter-like-api/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserUsecase struct {
	db             *gorm.DB
	log            *logrus.Logger
	config         *viper.Viper
	userRepository *repository.UserRepository
}

func NewUserUsecase(db *gorm.DB, log *logrus.Logger, config *viper.Viper, userRepository *repository.UserRepository) *UserUsecase {
	return &UserUsecase{
		db:             db,
		log:            log,
		config:         config,
		userRepository: userRepository,
	}
}

func (uc *UserUsecase) Register(ctx context.Context, request *model.UserRegisterRequest) (*model.UserResponse, error) {
	tx := uc.db.WithContext(ctx).Begin()
	defer tx.Rollback()

	total, err := uc.userRepository.CountById(tx, request.ID)
	if err != nil {
		uc.log.Warnf("error count user from database: %v", err)
		return nil, err
	}

	if total > 0 {
		uc.log.Warnf("user already exists: %v", err)
		return nil, echo.NewHTTPError(echo.ErrBadRequest.Code, fmt.Sprintf("user with id: %s already used", request.ID))
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		uc.log.Warnf("error hashing password: %v", err)
		return nil, echo.ErrInternalServerError
	}

	user := &entity.User{
		ID:       request.ID,
		Name:     request.Name,
		Email:    request.Email,
		Password: string(hashedPassword),
	}

	if err := uc.userRepository.Create(tx, user); err != nil {
		uc.log.Warnf("error create user: %v", err)
		return nil, echo.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		uc.log.Warnf("error commit database: %v", err)
		return nil, echo.ErrInternalServerError
	}

	return mapper.ToUserResponse(user), nil
}

func (uc *UserUsecase) Login(ctx context.Context, request *model.UserLoginRequest) (*model.UserTokenResponse, error) {
	tx := uc.db.WithContext(ctx).Begin()
	defer tx.Rollback()

	user := new(entity.User)
	if err := uc.userRepository.FindById(tx, user, request.ID); err != nil {
		uc.log.Warnf("error find user by id: %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, echo.NewHTTPError(echo.ErrUnauthorized.Code, fmt.Sprintf("user with id: %s not found", request.ID))
		}
		return nil, echo.ErrInternalServerError
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		uc.log.Warnf("error comparing password user: %v", err)
		return nil, echo.ErrUnauthorized
	}

	if err := tx.Commit().Error; err != nil {
		uc.log.Warnf("error commit database: %v", err)
		return nil, echo.ErrInternalServerError
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(12 * time.Hour).Unix(),
		"id":  user.ID,
	})

	secretKey := uc.config.GetString("jwt.secret")

	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		uc.log.Warnf("error signing token: %v", err)
		return nil, echo.ErrInternalServerError
	}

	return mapper.ToUserTokenResponse(signedToken), nil
}