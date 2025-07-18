package controllers

import (
	"net/http"
	"studyverse/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type RegisterInput struct {
	Username string `form:"username" binding:"required,min=3,max=20"`
	Email    string `form:"email" binding:"required,email"`
	Password string `form:"password" binding:"required,min=6"`
}

// GET: Kayıt sayfasını göster
func RegisterGet(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", nil)
}

// POST: Formdan gelen verileri al, kullanıcıyı oluştur
func RegisterPost(c *gin.Context) {
	var input RegisterInput

	// Formdan gelen veriyi kontrol et
	if err := c.ShouldBind(&input); err != nil {
		c.HTML(http.StatusBadRequest, "register.html", gin.H{
			"error": "Formda hata var. Lütfen geçerli bilgiler girin.",
		})
		return
	}

	// Şifreyi hashle
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "register.html", gin.H{
			"error": "Şifre işlenirken bir hata oluştu.",
		})
		return
	}

	// Kullanıcı nesnesi oluştur
	user := models.User{
		Username: input.Username,
		Email:    input.Email,
		Password: string(hashedPassword),
	}

	// Veritabanı bağlantısını al
	db, ok := c.MustGet("db").(*gorm.DB)
	if !ok {
		c.HTML(http.StatusInternalServerError, "register.html", gin.H{
			"error": "Veritabanı bağlantısı kurulamadı.",
		})
		return
	}

	// Kullanıcıyı veritabanına kaydet
	if err := db.Create(&user).Error; err != nil {
		c.HTML(http.StatusBadRequest, "register.html", gin.H{
			"error": "Kullanıcı oluşturulamadı. E-posta veya kullanıcı adı zaten kayıtlı olabilir.",
		})
		return
	}
}
