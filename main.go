package main

import (
	"log"
	"studyverse/controllers"
	"studyverse/models"
	"studyverse/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(".env dosyası yüklenemedi")
	}

	db, err := gorm.Open(sqlite.Open("studyverse.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Veritabanına bağlanılamadı:", err)
	}

	// Gerekli tabloları oluştur
	db.AutoMigrate(&models.User{}, &models.Task{}, &models.ResetToken{})

	r := gin.Default()

	// Veritabanını context'e ekle
	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	r.LoadHTMLGlob("templates/*")

	routes.SetupRoutes(r)

	authGroup := r.Group("/")
	authGroup.Use(controllers.AuthMiddleware())
	routes.RegisterTaskRoutes(authGroup)

	log.Println("Sunucu başlatıldı: http://localhost:6060")
	if err := r.Run(":6060"); err != nil {
		log.Fatal("Sunucu başlatılamadı:", err)
	}
}
