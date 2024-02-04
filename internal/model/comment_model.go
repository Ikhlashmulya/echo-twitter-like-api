package model

type CommentResponse struct {
	ID        string `json:"id"`
	User      string `json:"user"`
	Reply     string `json:"reply"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

type CommentCreateRequest struct {
	UserID string
	PostID string
	Reply  string `json:"reply" validate:"required"`
}

type CommentUpdateRequest struct {
	ID     string
	UserID string
	Reply  string `json:"reply" validate:"required"`
}

type CommentDeleteRequest struct {
	ID     string
	UserID string
}

type CommentFindByPostId struct {
	PostID string
	Page   int `query:"page"`
	Size   int `query:"size"`
}
