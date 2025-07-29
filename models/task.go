package models

import (
	"gorm.io/gorm"
)

type Task struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt string         `json:"created_at"`
	UpdatedAt string         `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	Date      string         `json:"date"`
	Content   string         `json:"content"`
	Done      bool           `json:"done"`
	UserID    uint           `json:"user_id"`
}

// CreateTask yeni bir task oluşturur
func CreateTask(db *gorm.DB, task *Task) error {
	return db.Create(task).Error
}

// GetTasksByUserAndDate kullanıcıya ve tarihe göre görevleri getirir
func GetTasksByUserAndDate(db *gorm.DB, userID uint, date string) ([]Task, error) {
	var tasks []Task
	err := db.Where("user_id = ? AND date = ?", userID, date).Order("done asc, created_at asc").Find(&tasks).Error
	return tasks, err
}

// UpdateTask güncelleme işlemi yapar
func UpdateTask(db *gorm.DB, task *Task) error {
	return db.Save(task).Error
}

// DeleteTask silme işlemi yapar
func DeleteTask(db *gorm.DB, taskID uint, userID uint) error {
	return db.Where("id = ? AND user_id = ?", taskID, userID).Delete(&Task{}).Error
}
