package model

type User struct {
	ID        string  `json:"id,omitempty" gorm:"column:id;primaryKey"`
	Email     string  `json:"email,omitempty" gorm:"column:email"`
	Password  string  `json:"password,omitempty" gorm:"column:password"`
	Token     string  `json:"token,omitempty" gorm:"column:-"`
	Followers []*User `json:"followers,omitempty" gorm:"many2many:user_followers"`
}
