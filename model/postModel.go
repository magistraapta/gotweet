package model

import "gorm.io/gorm"

type Post struct {
	gorm.Model
	ID       uint
	UserID   uint
	User     User      `gorm:"foreignKey:UserID"`
	Comments []Comment `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE;"`
	Content  string
}
