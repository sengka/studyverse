package controllers

import (
	"net/http"
	"strconv"
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

type LoginInput struct {
	Email    string `form:"email" binding:"required,email"`
	Password string `form:"password" binding:"required,min=6"`
}

// GET: Giriş sayfasını göster
func LoginGet(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

func LoginPost(c *gin.Context) {
	var input LoginInput

	if err := c.ShouldBind(&input); err != nil {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{
			"error": "Lütfen geçerli e-posta ve şifre girin.",
		})
		return
	}

	db, ok := c.MustGet("db").(*gorm.DB)
	if !ok {
		c.HTML(http.StatusInternalServerError, "login.html", gin.H{
			"error": "Veritabanı bağlantısı kurulamadı.",
		})
		return
	}

	var user models.User
	if err := db.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{
			"error": "E-posta bulunamadı.",
		})
		return
	}

	// Şifre doğrulama
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{
			"error": "Şifre yanlış.",
		})
		return
	}

	c.SetCookie("user_id", strconv.Itoa(int(user.ID)), 3600, "/", "localhost", false, true)

	c.Redirect(http.StatusSeeOther, "/homepage")
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := c.Cookie("user_id")
		if err != nil || userID == "" {
			// Cookie yoksa ya da boşsa kullanıcıyı login sayfasına gönder
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		// Cookie varsa devam et
		c.Next()
	}
}

func Logout(c *gin.Context) {
	// user_id çerezini sil
	c.SetCookie("user_id", "", -1, "/", "localhost", false, true)

	// Login sayfasına yönlendir
	c.Redirect(http.StatusSeeOther, "/greeting")
}
