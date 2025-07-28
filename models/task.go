package models

import "gorm.io/gorm"

type Task struct {
	gorm.Model
	Date    string `json:"date"`
	Content string `json:"content"`
	Done    bool   `json:"done"`
	UserID  uint   `json:"user_id"`
}
