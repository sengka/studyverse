package controllers

import (
	"log"
	"net/http"
	"studyverse/models"
	"time"

	"github.com/gin-gonic/gin"
)

func SaveTimelog(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		return
	}

	var input struct {
		Date     string `json:"date"`
		Duration int    `json:"duration"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("Geçersiz giriş: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz giriş"})
		return
	}

	timelog := models.Timelog{
		Date:     input.Date,
		Duration: input.Duration,
		UserID:   userID,
	}

	if err := models.DB.Create(&timelog).Error; err != nil {
		log.Printf("Süre kaydedilemedi: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Süre kaydedilemedi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Süre kaydedildi"})
}

func GetTimelogs(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		return
	}

	date := c.Query("date")
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	var timelogs []models.Timelog
	if err := models.DB.Where("user_id = ? AND date = ?", userID, date).Find(&timelogs).Error; err != nil {
		log.Printf("Kayıtlar alınamadı: user_id=%d, date=%s, Hata: %v", userID, date, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Kayıtlar alınamadı"})
		return
	}

	log.Printf("Dönen timelogs: %+v", timelogs)
	c.JSON(http.StatusOK, gin.H{"timelogs": timelogs})
}
