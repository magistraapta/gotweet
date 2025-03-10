package model

import "gorm.io/gorm"

type Comment struct {
	gorm.Model

	Content string `json:"content"`
	User    User   `gorm:"foreignKey:UserID"`
	Post    Post   `gorm:"foreignKey:PostID"`
	UserID  uint
	PostID  uint
}
