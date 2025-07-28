package models

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	database, err := gorm.Open(sqlite.Open("studyverse.db"), &gorm.Config{})
	if err != nil {
		panic("veritabanı bağlanamadı")
	}

	database.AutoMigrate(&User{}, &Task{}) // User modelin de burada
	DB = database
}
