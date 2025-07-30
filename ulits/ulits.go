package utils

import (
	"fmt"
	"log"
	"net/smtp"
	"os"

	"github.com/joho/godotenv"
)

func SendEmail(to, subject, body string) error {
	// .env dosyasını yükle (eğer zaten ana programda yüklüyorsan buraya gerek yok)
	err := godotenv.Load()
	if err != nil {
		log.Println("Env yüklenemedi:", err)
	}

	// SMTP ayarlarını env'den al
	from := os.Getenv("EMAIL_ADDRESS")
	password := os.Getenv("EMAIL_PASSWORD")

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Mesaj formatı
	message := []byte(fmt.Sprintf(
		"From: %s\r\n"+
			"To: %s\r\n"+
			"Subject: %s\r\n"+
			"\r\n%s\r\n", from, to, subject, body))

	// SMTP auth oluştur
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Mail gönder
	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, message)
	if err != nil {
		return err
	}
	return nil
}
