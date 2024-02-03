package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Comment struct {
	ID        string `gorm:"column:id;primaryKey"`
	UserID    string `gorm:"column:user_id"`
	PostID    string `gorm:"column:post_id"`
	Reply     string `gorm:"column:reply"`
	CreatedAt int64  `gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt int64  `gorm:"column:updated_at;autoCreateTime:milli;autoUpdateTime:milli"`
	Post      Post   `gorm:"foreignKey:post_id;references:id"`
	User      User   `gorm:"foreignKey:user_id;references:id"`
}

func (c *Comment) BeforeCreate(tx *gorm.DB) (err error) {
	if c.ID == "" {
		c.ID = uuid.NewString()
	}

	return
}
