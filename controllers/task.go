package controllers

import (
	"net/http"
	"strconv"
	"studyverse/models"

	"github.com/gin-gonic/gin"
)

func CalendarPage(c *gin.Context) {
	c.HTML(http.StatusOK, "calendar.html", nil)
}

func GetTasks(c *gin.Context) {
	date := c.Query("date")
	userIDStr, _ := c.Cookie("user_id")
	userID, _ := strconv.Atoi(userIDStr)

	var tasks []models.Task
	models.DB.Where("date = ? AND user_id = ?", date, userID).Find(&tasks)
	c.JSON(http.StatusOK, tasks)
}

func CreateTask(c *gin.Context) {
	var input struct {
		Date    string `json:"date"`
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Veri hatalı"})
		return
	}

	userIDStr, _ := c.Cookie("user_id")
	userID, _ := strconv.Atoi(userIDStr)

	task := models.Task{
		Date:    input.Date,
		Content: input.Content,
		UserID:  uint(userID),
		Done:    false,
	}
	models.DB.Create(&task)
	c.JSON(http.StatusOK, task)
}

func DeleteTask(c *gin.Context) {
	id := c.Param("id")
	models.DB.Delete(&models.Task{}, id)
	c.JSON(http.StatusOK, gin.H{"message": "Silindi"})
}

func UpdateTask(c *gin.Context) {
	id := c.Param("id")

	var input struct {
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Veri hatalı"})
		return
	}

	var task models.Task
	if err := models.DB.First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Görev bulunamadı"})
		return
	}

	task.Content = input.Content
	models.DB.Save(&task)

	c.JSON(http.StatusOK, task)
}

func ToggleDone(c *gin.Context) {
	id := c.Param("id")
	var task models.Task
	if err := models.DB.First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Görev bulunamadı"})
		return
	}

	task.Done = !task.Done
	models.DB.Save(&task)

	c.JSON(http.StatusOK, task)
}
