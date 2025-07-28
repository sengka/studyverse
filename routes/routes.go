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

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "greeting.html", nil)
	})
	r.GET("/explore", func(c *gin.Context) {
		c.HTML(http.StatusOK, "explore.html", nil)
	})
	r.GET("/forgotpassword", func(c *gin.Context) {
		c.HTML(http.StatusOK, "forgotpassword.html", nil)
	})
	r.GET("/login", controllers.LoginGet)
	r.POST("/login", controllers.LoginPost)
	r.GET("/register", controllers.RegisterGet)
	r.POST("/register", controllers.RegisterPost)

	// Korunan rotalar
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
	auth.GET("/calendar", controllers.CalendarPage)

}

func RegisterTaskRoutes(rg *gin.RouterGroup) {
	tasks := rg.Group("/api/tasks")
	tasks.GET("", controllers.GetTasks)
	tasks.POST("", controllers.CreateTask)
	tasks.DELETE("/:id", controllers.DeleteTask)
	tasks.PUT("/:id", controllers.UpdateTask)
	tasks.PATCH("/:id/done", controllers.ToggleDone)
}
