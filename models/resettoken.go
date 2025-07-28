package models

import (
	"time"
)

type ResetToken struct {
	Token     string    `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null"`
	ExpiresAt time.Time `gorm:"not null"`
}
