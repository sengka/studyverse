package main

import (
	"log"
	"studyverse/models"
	"studyverse/routes"
	"studyverse/services" // burada servisleri ekliyoruz

	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(".env dosyası yüklenemedi:", err)
	}

	log.Println("Veritabanı başlatılıyor...")
	models.InitDB()
	log.Println("Veritabanı başarıyla başlatıldı")

	// Günlük görev hatırlatma e-postalarını arka planda başlatıyoruz
	go func() {
		for {
			services.CheckAndSendDailyTasks(models.DB)
			time.Sleep(24 * time.Hour)
		}
	}()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:6060"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.Use(func(c *gin.Context) {
		log.Printf("İstek alındı: %s %s", c.Request.Method, c.Request.URL)
		c.Set("db", models.DB)
		c.Next()
	})

	r.LoadHTMLGlob("templates/*")
	log.Println("Rotalar yükleniyor...")
	routes.SetupRoutes(r)
	log.Println("Rotalar yüklendi")

	log.Println("Sunucu başlatıldı: http://localhost:6060")
	if err := r.Run(":6060"); err != nil {
		log.Fatal("Sunucu başlatılamadı:", err)
	}
}
