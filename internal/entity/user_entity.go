package entity

type User struct {
	ID           string     `gorm:"column:id;primaryKey;size:255"`
	Name         string     `gorm:"column:name"`
	Email        string     `gorm:"column:email"`
	Password     string     `gorm:"column:password"`
	PhotoProfile string     `gorm:"column:photo_profile"`
	Followers    []*User    `gorm:"many2many:user_followers"`
	Posts        []*Post    `gorm:"foreignKey:user_id;references:id"`
	Likes        []*Post    `gorm:"many2many:user_likes_post;foreignKey:id;joinForeignKey:user_id;references:id;joinReferences:post_id"`
	Comments     []*Comment `gorm:"foreignKey:user_id;references:id"`
}
