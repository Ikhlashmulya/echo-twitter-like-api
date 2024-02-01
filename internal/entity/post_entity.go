package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Post struct {
	ID        string     `gorm:"column:id;primaryKey;size:255"`
	UserID    string     `gorm:"column:user_id"`
	Content   string     `gorm:"column:content"`
	CreatedAt int64      `gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt int64      `gorm:"column:updated_at;autoCreateTime:milli;autoUpdateTime:milli"`
	Author    User       `gorm:"foreignKey:user_id;references:id"`
	Liked     []*User    `gorm:"many2many:user_likes_post;foreignKey:id;joinForeignKey:user_id;references:id;joinReferences:post_id"`
	Comments  []*Comment `gorm:"foreignKey:post_id;references:id"`
}

func (p *Post) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == "" {
		p.ID = uuid.NewString()
	}
	
	return
}
