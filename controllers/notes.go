package controllers

import (
	"net/http"
	"strconv"

	"studyverse/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetNotes tüm notları ve klasörleri döndürür
func GetNotes(c *gin.Context) {
	db, _ := c.MustGet("db").(*gorm.DB)
	userID, _ := c.Cookie("user_id")
	uid64, _ := strconv.ParseUint(userID, 10, 32)

	var folders []models.Folder
	var notes []models.Note
	if err := db.Where("user_id = ?", uint(uid64)).Find(&folders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := db.Where("user_id = ?", uint(uid64)).Find(&notes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{Folders: folders, Notes: notes})
}

// CreateFolder yeni bir klasör ekler
func CreateFolder(c *gin.Context) {
	db, _ := c.MustGet("db").(*gorm.DB)
	userID, _ := c.Cookie("user_id")
	uid64, _ := strconv.ParseUint(userID, 10, 32)

	var folder models.Folder
	if err := c.ShouldBindJSON(&folder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	folder.UserID = uint(uid64)

	if err := db.Create(&folder).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, folder)
}

func CreateNote(c *gin.Context) {
	db, _ := c.MustGet("db").(*gorm.DB)
	userID, _ := c.Cookie("user_id")
	uid64, _ := strconv.ParseUint(userID, 10, 32)

	var note models.Note
	if err := c.ShouldBindJSON(&note); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	note.UserID = uint(uid64)

	// folder_id string olarak geldiyse tam sayıya çevir
	if note.FolderID != nil {
		folderIDStr := strconv.FormatUint(uint64(*note.FolderID), 10)
		folderID, err := strconv.ParseUint(folderIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid folder_id"})
			return
		}
		*note.FolderID = uint(folderID)
	}

	if err := db.Create(&note).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, note)
}

// UpdateFolder mevcut bir klasörü günceller
func UpdateFolder(c *gin.Context) {
	db, _ := c.MustGet("db").(*gorm.DB)
	userID, _ := c.Cookie("user_id")
	uid64, _ := strconv.ParseUint(userID, 10, 32)
	id := c.Param("id")
	folderID, _ := strconv.Atoi(id)

	var folder models.Folder
	if err := db.Where("id = ? AND user_id = ?", folderID, uint(uid64)).First(&folder).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Folder not found"})
		return
	}

	if err := c.ShouldBindJSON(&folder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	folder.UserID = uint(uid64)

	if err := db.Save(&folder).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, folder)
}

// UpdateNote mevcut bir notu günceller
func UpdateNote(c *gin.Context) {
	db, _ := c.MustGet("db").(*gorm.DB)
	userID, _ := c.Cookie("user_id")
	uid64, _ := strconv.ParseUint(userID, 10, 32)
	id := c.Param("id")
	noteID, _ := strconv.Atoi(id)

	var note models.Note
	if err := db.Where("id = ? AND user_id = ?", noteID, uint(uid64)).First(&note).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
		return
	}

	if err := c.ShouldBindJSON(&note); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	note.UserID = uint(uid64)

	if err := db.Save(&note).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, note)
}

// DeleteFolder klasörü siler
func DeleteFolder(c *gin.Context) {
	db, _ := c.MustGet("db").(*gorm.DB)
	userID, _ := c.Cookie("user_id")
	uid64, _ := strconv.ParseUint(userID, 10, 32)
	id := c.Param("id")
	folderID, _ := strconv.Atoi(id)

	var folder models.Folder
	if err := db.Where("id = ? AND user_id = ?", folderID, uint(uid64)).First(&folder).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Folder not found"})
		return
	}

	if err := db.Delete(&folder).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// DeleteNote notu siler
func DeleteNote(c *gin.Context) {
	db, _ := c.MustGet("db").(*gorm.DB)
	userID, _ := c.Cookie("user_id")
	uid64, _ := strconv.ParseUint(userID, 10, 32)
	id := c.Param("id")
	noteID, _ := strconv.Atoi(id)

	var note models.Note
	if err := db.Where("id = ? AND user_id = ?", noteID, uint(uid64)).First(&note).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
		return
	}

	if err := db.Delete(&note).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
