package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"size:100;not null"`
	Email    string `gorm:"uniqueIndex;not null"`
	Password string `gorm:"olumn:password;not null"`
}
