package entity

type User struct {
	ID           string `gorm:"column:id;primaryKey"`
	Name         string `gorm:"column:name"`
	Email        string `gorm:"column:email"`
	Password     string `gorm:"column:password"`
	PhotoProfile string `gorm:"column:photo_profile"`
}
