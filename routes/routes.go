package routes

import (
	"net/http"
	"strconv"

	"studyverse/controllers"
	"studyverse/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine) {
	r.LoadHTMLGlob("templates/*")

	// Genel Sayfalar
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "greeting.html", nil)
	})
	r.GET("/explore", func(c *gin.Context) {
		c.HTML(http.StatusOK, "explore.html", nil)
	})

	// Auth Sayfaları
	r.GET("/login", controllers.LoginGet)
	r.POST("/login", controllers.LoginPost)
	r.GET("/register", controllers.RegisterGet)
	r.POST("/register", controllers.RegisterPost)

	// Şifremi unuttum
	r.GET("/forgotpassword", func(c *gin.Context) {
		c.HTML(http.StatusOK, "forgotpassword.html", nil)
	})
	r.POST("/forgotpassword", controllers.ForgotPasswordPost)
	r.GET("/changepassword", controllers.NewPasswordGet)
	r.POST("/changepassword", controllers.NewPasswordPost)

	// Giriş gerektiren (korumalı) sayfalar için grup
	auth := r.Group("/")
	auth.Use(controllers.AuthMiddleware())

	auth.GET("/homepage", func(c *gin.Context) {
		db, _ := c.MustGet("db").(*gorm.DB)
		userID, _ := c.Cookie("user_id")
		uid64, _ := strconv.ParseUint(userID, 10, 32)

		var user models.User
		if err := db.First(&user, uint(uid64)).Error; err != nil {
			c.Redirect(http.StatusFound, "/login")
			return
		}

		c.HTML(http.StatusOK, "homepage.html", gin.H{
			"Username": user.Username,
		})
	})

	auth.GET("/logout", controllers.Logout)

	// Calendar sayfası
	auth.GET("/calendar", func(c *gin.Context) {
		db, _ := c.MustGet("db").(*gorm.DB)
		userID, _ := c.Cookie("user_id")
		uid64, _ := strconv.ParseUint(userID, 10, 32)

		var user models.User
		if err := db.First(&user, uint(uid64)).Error; err != nil {
			c.Redirect(http.StatusFound, "/login")
			return
		}

		c.HTML(http.StatusOK, "calendar.html", gin.H{
			"Username": user.Username,
		})
	})

	// Todo sayfası
	auth.GET("/todo", func(c *gin.Context) {
		db, _ := c.MustGet("db").(*gorm.DB)
		userID, _ := c.Cookie("user_id")
		uid64, _ := strconv.ParseUint(userID, 10, 32)

		var user models.User
		if err := db.First(&user, uint(uid64)).Error; err != nil {
			c.Redirect(http.StatusFound, "/login")
			return
		}

		c.HTML(http.StatusOK, "todo.html", gin.H{
			"Username": user.Username,
		})
	})

	// Timer sayfası
	auth.GET("/timer", func(c *gin.Context) {
		db, _ := c.MustGet("db").(*gorm.DB)
		userID, _ := c.Cookie("user_id")
		uid64, _ := strconv.ParseUint(userID, 10, 32)

		var user models.User
		if err := db.First(&user, uint(uid64)).Error; err != nil {
			c.Redirect(http.StatusFound, "/login")
			return
		}

		c.HTML(http.StatusOK, "timer.html", gin.H{
			"Username": user.Username,
		})
	})

	// Notes sayfası
	auth.GET("/notes", func(c *gin.Context) {
		db, _ := c.MustGet("db").(*gorm.DB)
		userID, _ := c.Cookie("user_id")
		uid64, _ := strconv.ParseUint(userID, 10, 32)

		var user models.User
		if err := db.First(&user, uint(uid64)).Error; err != nil {
			c.Redirect(http.StatusFound, "/login")
			return
		}

		c.HTML(http.StatusOK, "notes.html", gin.H{
			"Username": user.Username,
		})
	})

	// API rotaları için alt grup (prefix: /api/tasks)
	taskGroup := auth.Group("/api/tasks")
	taskGroup.GET("", controllers.GetTasks)          // GET /api/tasks?date=yyyy-mm-dd
	taskGroup.POST("", controllers.CreateTask)       // POST /api/tasks
	taskGroup.PUT("/:id", controllers.UpdateTask)    // PUT /api/tasks/:id
	taskGroup.DELETE("/:id", controllers.DeleteTask) // DELETE /api/tasks/:id

	// Timer API rotaları
	r.POST("/api/timelogs", controllers.SaveTimelog)
	r.GET("/api/timelogs", controllers.GetTimelogs)

	// Event API rotaları
	r.POST("/api/events", controllers.SaveEvent)
	r.GET("/api/events", controllers.GetEvents)
	r.DELETE("/api/events/:id", controllers.DeleteEvent)

	// Notes API rotaları
	notesGroup := auth.Group("/api/notes")
	notesGroup.GET("", controllers.GetNotes)
	notesGroup.POST("/folders", controllers.CreateFolder)
	notesGroup.POST("", controllers.CreateNote)
	notesGroup.PUT("/folders/:id", controllers.UpdateFolder)
	notesGroup.PUT("/:id", controllers.UpdateNote)
	notesGroup.DELETE("/folders/:id", controllers.DeleteFolder)
	notesGroup.DELETE("/:id", controllers.DeleteNote)
}
