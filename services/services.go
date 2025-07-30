package services

import (
	"fmt"
	"log"
	"time"

	"studyverse/models"
	utils "studyverse/ulits"

	"gorm.io/gorm"
)

func CheckAndSendDailyTasks(db *gorm.DB) {
	var users []models.User
	db.Find(&users)

	today := time.Now().Format("2006-01-02")

	for _, user := range users {
		var tasks []models.Task // Task modelinde DueDate alanı olmalı
		db.Where("user_id = ? AND date = ?", user.ID, today).Find(&tasks)

		if len(tasks) > 0 {
			body := "Bugün yapman gereken görevler:\n\n"
			for _, task := range tasks {
				body += fmt.Sprintf("- %s\n", task.Content)
			}

			err := utils.SendEmail(user.Email, "StudyVerse - Bugünkü Görevlerin", body)
			if err != nil {
				log.Println("E-posta gönderilemedi:", err)
			} else {
				log.Println("E-posta gönderildi:", user.Email)
			}
		}
	}
}
