package main

import (
	"log"
	"studyverse/controllers"
	"studyverse/models"
	"studyverse/routes"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// Veritabanı bağlantısı ve migrate
	db, err := gorm.Open(sqlite.Open("studyverse.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Veritabanına bağlanılamadı:", err)
	}

	// Gerekli tabloları oluştur
	db.AutoMigrate(&models.User{}, &models.Task{})

	// Eğer models.InitDB fonksiyonun veritabanını başlatmak içinse,
	// onu burada db ile parametre olarak çağırmalısın:
	models.InitDB()

	r := gin.Default()

	// Veritabanını context'e ekle (her request'te erişilebilir olsun diye)
	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	// HTML şablonlarını yükle
	r.LoadHTMLGlob("templates/*")

	// Tüm route'ları burada organize et
	routes.SetupRoutes(r)

	// Task ile ilgili route'ları auth middleware ile koru
	authGroup := r.Group("/")
	authGroup.Use(controllers.AuthMiddleware())
	routes.RegisterTaskRoutes(authGroup)

	// Sunucuyu başlat
	log.Println("Sunucu başlatıldı: http://localhost:6060")
	if err := r.Run(":6060"); err != nil {
		log.Fatal("Sunucu başlatılamadı:", err)
	}
}
