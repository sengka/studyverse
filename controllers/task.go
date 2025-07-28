package controllers

import (
	"net/http"
	"strconv"
	"time"

	"studyverse/models"

	"github.com/gin-gonic/gin"
)

func getUserIDFromContext(c *gin.Context) (uint, bool) {
	userIDStr, err := c.Cookie("user_id")
	if err != nil || userIDStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Kullanıcı ID'si bulunamadı"})
		return 0, false
	}

	userID64, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz kullanıcı ID'si"})
		return 0, false
	}
	return uint(userID64), true
}

// GET /tasks?date=yyyy-mm-dd
func GetTasks(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		return
	}

	date := c.Query("date")
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	tasks, err := models.GetTasksByUserAndDate(models.DB, userID, date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Görevler alınamadı"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tasks": tasks})
}

// POST /tasks
func CreateTask(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		return
	}

	var input struct {
		Date    string `json:"date" binding:"required"`
		Content string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz veri"})
		return
	}

	task := models.Task{
		Date:    input.Date,
		Content: input.Content,
		Done:    false,
		UserID:  userID,
	}

	if err := models.CreateTask(models.DB, &task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Görev oluşturulamadı"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"task": task})
}

// PUT /tasks/:id
func UpdateTask(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		return
	}

	taskID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz görev ID'si"})
		return
	}

	var input struct {
		Content *string `json:"content"`
		Done    *bool   `json:"done"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz veri"})
		return
	}

	// Görevi çek
	var task models.Task
	err = models.DB.Where("id = ? AND user_id = ?", uint(taskID), userID).First(&task).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Görev bulunamadı"})
		return
	}

	if input.Content != nil {
		task.Content = *input.Content
	}
	if input.Done != nil {
		task.Done = *input.Done
	}

	if err := models.UpdateTask(models.DB, &task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Görev güncellenemedi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"task": task})
}

// DELETE /tasks/:id
func DeleteTask(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		return
	}

	taskID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz görev ID'si"})
		return
	}

	err = models.DeleteTask(models.DB, uint(taskID), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Görev silinemedi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Görev silindi"})
}
