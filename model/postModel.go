package model

import "gorm.io/gorm"

type Post struct {
	gorm.Model
	UserID   uint
	User     User      `gorm:"foreignKey:UserID"`
	Comments []Comment `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE;"`
	Content  string
}
