package models

import (
	"gorm.io/gorm"
)

type Timelog struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Date      string         `json:"date"`
	Duration  int            `json:"duration"`
	UserID    uint           `json:"user_id"`
	CreatedAt string         `json:"created_at"`
	UpdatedAt string         `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
