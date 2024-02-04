package model

type PostCreateRequest struct {
	UserID  string `json:"-"`
	Content string `json:"content" validate:"required"`
}

type PostUpdateRequest struct {
	ID      string `json:"-"`
	UserID  string `json:"-"`
	Content string `json:"content" validate:"required"`
}

type PostDeleteRequest struct {
	ID     string
	UserID string
}

type PostFindByUserIdRequest struct {
	UserID string 
	Page   int `query:"page"`
	Size   int `query:"size"`
}

type PostFindByFollowingUserRequest struct {
	UserID string
	Page   int `query:"page"`
	Size   int `query:"size"`
}

type PostResponse struct {
	ID        string `json:"id"`
	UserID    string `json:"user"`
	Content   string `json:"content"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}
