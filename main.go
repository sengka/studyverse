package main

import (
	"log"
	"net/http"
	"strconv"
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

	// Giriş yapılması gerekmeyen sayfalar
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "greeting.html", nil)
	})
	r.GET("/explore", func(c *gin.Context) { // <-- Artık burada, korumasız
		c.HTML(http.StatusOK, "explore.html", nil)
	})
	r.GET("/login", controllers.LoginGet)
	r.POST("/login", controllers.LoginPost)
	r.GET("/register", controllers.RegisterGet)
	r.POST("/register", controllers.RegisterPost)

	// Giriş yapılması gereken (korumalı) sayfalar
	authRoutes := r.Group("/")
	authRoutes.Use(controllers.AuthMiddleware())
	{
		authRoutes.GET("/homepage", func(c *gin.Context) {
			db, _ := c.MustGet("db").(*gorm.DB)

			// Cookie'den user_id'yi al
			userID, err := c.Cookie("user_id")
			if err != nil {
				c.Redirect(http.StatusFound, "/login")
				return
			}

			// user_id string, onu uint'a çevir
			uid64, _ := strconv.ParseUint(userID, 10, 32)
			var user models.User
			if err := db.First(&user, uint(uid64)).Error; err != nil {
				c.Redirect(http.StatusFound, "/login")
				return
			}

			// HTML'e user bilgisini gönder
			c.HTML(http.StatusOK, "homepage.html", gin.H{
				"Username": user.Username, // Veritabanındaki kullanıcı adı
			})
		})
		r.GET("/logout", controllers.Logout)

	}

	log.Println("Sunucu http://localhost:6060 adresinde başlatıldı. Ctrl+C ile durdurabilirsiniz.")

	if err := r.Run(":6060"); err != nil {
		log.Fatal("Sunucu başlatılamadı:", err)
	}
}
