package main

import (
	"log"
	"net/http"
	"studyverse/controllers"
	"studyverse/models"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// Veritabanı bağlan ve migrate et
	db, err := gorm.Open(sqlite.Open("studyverse.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Veritabanına bağlanılamadı:", err)
	}
	db.AutoMigrate(&models.User{})

	r := gin.Default()

	// HTML şablonlarını yükle
	r.LoadHTMLGlob("templates/*")

	// Veritabanını context'e ekle
	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	// Routes
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "greeting.html", nil)
	})

	r.GET("/register", controllers.RegisterGet)
	r.POST("/register", controllers.RegisterPost)

	r.GET("/login", controllers.LoginGet)
	r.POST("/login", controllers.LoginPost)

	log.Println("Sunucu http://localhost:6060 adresinde başlatıldı. Ctrl+C ile durdurabilirsiniz.")

	if err := r.Run(":6060"); err != nil {
		log.Fatal("Sunucu başlatılamadı:", err)
	}
}
