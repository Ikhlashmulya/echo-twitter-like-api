package model

type Post struct {
	ID      string `json:"id,omitempty" gorm:"column:id;primaryKey"`
	To      string `json:"to,omitempty"gorm:"column:to"`
	From    string `json:"from,omitempty"gorm:"column:from"`
	Message string `json:"message,omitempty"gorm:"column:message"`
}
