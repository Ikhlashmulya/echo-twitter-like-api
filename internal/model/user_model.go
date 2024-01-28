package model

type UserRegisterRequest struct {
	ID       string `json:"id" validate:"required"`
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UserLoginRequest struct {
	ID       string `json:"id" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type FindAllFollowerRequest struct {
	UserID string `json:"user_id"`
	Page   int    `json:"-" query:"page"`
	Size   int    `json:"-" query:"size"`
}

type UserResponse struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	PhotoProfile string `json:"photo_profile"`
}

type UserTokenResponse struct {
	Token string `json:"token"`
}
