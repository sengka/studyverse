package controllers

import (
	"log"
	"net/http"
	"studyverse/models"

	"github.com/gin-gonic/gin"
)

func SaveEvent(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		log.Printf("Kullanıcı ID'si alınamadı")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Kullanıcı ID'si bulunamadı"})
		return
	}

	var input struct {
		Date  string `json:"date"`
		Title string `json:"title"`
		Type  string `json:"type"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("Geçersiz giriş: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz giriş verisi"})
		return
	}

	if input.Date == "" || input.Title == "" || input.Type == "" {
		log.Printf("Eksik alanlar: date=%s, title=%s, type=%s", input.Date, input.Title, input.Type)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tarih, başlık ve tür zorunludur"})
		return
	}

	// Geçerli etkinlik türlerini kontrol et
	validTypes := map[string]bool{"exam": true, "meeting": true, "assignment": true, "project": true, "other": true}
	if !validTypes[input.Type] {
		log.Printf("Geçersiz etkinlik türü: %s", input.Type)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz etkinlik türü"})
		return
	}

	event := models.Event{
		Date:   input.Date,
		Title:  input.Title,
		Type:   input.Type,
		UserID: userID,
	}

	if err := models.DB.Create(&event).Error; err != nil {
		log.Printf("Etkinlik kaydedilemedi: user_id=%d, date=%s, title=%s, type=%s, Hata: %v", userID, input.Date, input.Title, input.Type, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Etkinlik kaydedilemedi"})
		return
	}

	log.Printf("Etkinlik kaydedildi: user_id=%d, date=%s, title=%s, type=%s", userID, input.Date, input.Title, input.Type)
	c.JSON(http.StatusOK, gin.H{"message": "Etkinlik kaydedildi", "event": event})
}

func GetEvents(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		log.Printf("Kullanıcı ID'si alınamadı")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Kullanıcı ID'si bulunamadı"})
		return
	}

	date := c.Query("date")
	var events []models.Event
	query := models.DB.Where("user_id = ?", userID)
	if date != "" {
		query = query.Where("date = ?", date)
	}
	if err := query.Find(&events).Error; err != nil {
		log.Printf("Etkinlikler alınamadı: user_id=%d, date=%s, Hata: %v", userID, date, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Etkinlikler alınamadı"})
		return
	}

	log.Printf("Dönen etkinlikler: user_id=%d, count=%d", userID, len(events))
	c.JSON(http.StatusOK, gin.H{"events": events})
}

func DeleteEvent(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		log.Printf("Kullanıcı ID'si alınamadı")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Kullanıcı ID'si bulunamadı"})
		return
	}

	id := c.Param("id")
	var event models.Event
	if err := models.DB.Where("id = ? AND user_id = ?", id, userID).First(&event).Error; err != nil {
		log.Printf("Etkinlik bulunamadı: id=%s, user_id=%d, Hata: %v", id, userID, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Etkinlik bulunamadı"})
		return
	}

	if err := models.DB.Delete(&event).Error; err != nil {
		log.Printf("Etkinlik silinemedi: id=%s, user_id=%d, Hata: %v", id, userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Etkinlik silinemedi"})
		return
	}

	log.Printf("Etkinlik silindi: id=%s, user_id=%d", id, userID)
	c.JSON(http.StatusOK, gin.H{"message": "Etkinlik silindi"})
}
