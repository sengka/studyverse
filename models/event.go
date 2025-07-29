package models

import (
	"time"

	"gorm.io/gorm"
)

type Event struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Date      string         `json:"date"`
	Title     string         `json:"title"`
	Type      string         `json:"type"`
	UserID    uint           `json:"user_id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
