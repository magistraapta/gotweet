package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID       uint `gorm:"primary_key"`
	Username string
	Password string
	Email    string `gorm:"unique"`
}
