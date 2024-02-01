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

type PostResponse struct {
	ID        string `json:"id"`
	UserID    string `json:"author"`
	Content   string `json:"content"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}
