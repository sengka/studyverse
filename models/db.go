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
		DisableForeignKeyConstraintWhenMigrating: true, // SQLite için yabancı anahtar kısıtlamalarını devre dışı bırak
	})
	if err != nil {
		log.Fatalf("Veritabanına bağlanılamadı: %v", err)
	}

	// Bağlantı havuzu ayarları (isteğe bağlı)
	sqlDB, err := database.DB()
	if err != nil {
		log.Fatalf("Veritabanı bağlantı havuzu alınamadı: %v", err)
	}
	sqlDB.SetMaxOpenConns(10) // Maksimum açık bağlantı
	sqlDB.SetMaxIdleConns(5)  // Maksimum boşta bağlantı

	// Modelleri migre et
	err = database.AutoMigrate(&User{}, &Task{}, &ResetToken{}, &Timelog{})
	if err != nil {
		log.Fatalf("Tablolar migre edilemedi: %v", err)
	}

	DB = database
	log.Println("Veritabanı başarıyla başlatıldı")
}
