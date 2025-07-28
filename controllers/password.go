package controllers

import (
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"time"

	"studyverse/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Şifre sıfırlama için e-posta gönderme
func ForgotPasswordPost(c *gin.Context) {
	db, _ := c.MustGet("db").(*gorm.DB)
	email := c.PostForm("email")

	var user models.User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		c.HTML(http.StatusOK, "forgotpassword.html", gin.H{
			"error": "Bu e-posta adresine ait bir kullanıcı bulunamadı.",
		})
		return
	}

	token := uuid.NewString()
	resetToken := models.ResetToken{
		Token:     token,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(15 * time.Minute),
	}

	if err := db.Create(&resetToken).Error; err != nil {
		log.Println("Token kaydedilemedi:", err)
		c.HTML(http.StatusOK, "forgotpassword.html", gin.H{
			"error": "Sunucu hatası. Lütfen tekrar deneyin.",
		})
		return
	}

	resetLink := fmt.Sprintf("http://localhost:6060/changepassword?token=%s", token)
	subject := "Şifre Sıfırlama Talebi"
	body := fmt.Sprintf("Merhaba %s,\n\nŞifrenizi sıfırlamak için aşağıdaki bağlantıya tıklayın:\n\n%s\n\nBu bağlantı 15 dakika içinde geçerliliğini yitirecektir.\n\nİyi çalışmalar! 🚀", user.Username, resetLink)

	message := []byte("Subject: " + subject + "\r\n\r\n" + body)

	auth := smtp.PlainAuth("", os.Getenv("EMAIL_ADDRESS"), os.Getenv("EMAIL_PASSWORD"), "smtp.gmail.com")
	err := smtp.SendMail("smtp.gmail.com:587", auth, os.Getenv("EMAIL_ADDRESS"), []string{user.Email}, message)
	if err != nil {
		log.Println("SMTP error:", err)
		c.HTML(http.StatusOK, "forgotpassword.html", gin.H{
			"error": "E-posta gönderilirken bir hata oluştu.",
		})
		return
	}

	c.HTML(http.StatusOK, "forgotpassword.html", gin.H{
		"success": "Şifre sıfırlama bağlantısı e-posta adresinize gönderildi!",
	})
}

// Şifre sıfırlama sayfasını göster
func NewPasswordGet(c *gin.Context) {
	db, _ := c.MustGet("db").(*gorm.DB)
	token := c.Query("token")

	if token == "" {
		c.String(http.StatusBadRequest, "Token gerekli.")
		return
	}

	var resetToken models.ResetToken
	if err := db.Where("token = ?", token).First(&resetToken).Error; err != nil || resetToken.ExpiresAt.Before(time.Now()) {
		c.HTML(http.StatusOK, "changepassword.html", gin.H{
			"error": "Geçersiz veya süresi dolmuş token.",
		})
		return
	}

	c.HTML(http.StatusOK, "changepassword.html", gin.H{
		"Token": token,
	})
}

// Yeni şifreyi alıp veritabanında güncelle
func NewPasswordPost(c *gin.Context) {
	db, _ := c.MustGet("db").(*gorm.DB)
	token := c.PostForm("token")
	password := c.PostForm("password")

	if len(password) < 6 {
		c.HTML(http.StatusOK, "changepassword.html", gin.H{
			"token": token,
			"error": "Şifre en az 6 karakter olmalı.",
		})
		return
	}

	var resetToken models.ResetToken
	if err := db.Where("token = ?", token).First(&resetToken).Error; err != nil || resetToken.ExpiresAt.Before(time.Now()) {
		c.HTML(http.StatusOK, "changepassword.html", gin.H{
			"error": "Geçersiz veya süresi dolmuş token.",
		})
		return
	}

	// Şifreyi hash'le
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Şifre hashleme hatası:", err)
		c.HTML(http.StatusOK, "changepassword.html", gin.H{
			"error": "Sunucu hatası. Tekrar deneyin.",
			"token": token,
		})
		return
	}

	// Şifreyi güncelle
	if err := db.Model(&models.User{}).Where("id = ?", resetToken.UserID).Update("password", string(hashedPassword)).Error; err != nil {
		log.Println("Şifre güncelleme hatası:", err)
		c.HTML(http.StatusOK, "changepassword.html", gin.H{
			"error": "Şifre güncellenemedi.",
			"token": token,
		})
		return
	}
	fmt.Println("Yeni şifre veritabanına kaydedildi")

	// Token'ı geçersiz kıl
	db.Delete(&models.ResetToken{}, "token = ?", token)

	c.HTML(http.StatusOK, "changepassword.html", gin.H{
		"success": "Şifreniz başarıyla güncellendi. Giriş yapabilirsiniz!",
	})

}
