package models

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	database, err := gorm.Open(sqlite.Open("studyverse.db"), &gorm.Config{
		// GORM ayarları
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		log.Fatalf("Veritabanına bağlanılamadı: %v", err)
	}

	// Bağlantı havuzu ayarları (isteğe bağlı)
	sqlDB, err := database.DB()
	if err != nil {
		log.Fatalf("Veritabanı bağlantı havuzu alınamadı: %v", err)
	}
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)

	// Modelleri migre et
	err = database.AutoMigrate(&User{}, &Task{}, &ResetToken{}, &Timelog{}, &Event{})
	if err != nil {
		log.Fatalf("Tablolar migre edilemedi: %v", err)
	}

	DB = database
	log.Println("Veritabanı başarıyla başlatıldı")
}
