package main

import (
	"log"
	"studyverse/models"
	"studyverse/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(".env dosyası yüklenemedi")
	}

	// Veritabanını başlat (models.DB set edilir ve AutoMigrate burada yapılır)
	models.InitDB()

	r := gin.Default()

	// Global models.DB'yi context'e ekle
	r.Use(func(c *gin.Context) {
		c.Set("db", models.DB)
		c.Next()
	})

	r.LoadHTMLGlob("templates/*")

	routes.SetupRoutes(r)

	log.Println("Sunucu başlatıldı: http://localhost:6060")
	if err := r.Run(":6060"); err != nil {
		log.Fatal("Sunucu başlatılamadı:", err)
	}
}
