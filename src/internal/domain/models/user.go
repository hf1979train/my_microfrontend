package models

import (
	"time"
)

type User struct {
	UserID    uint      `gorm:"primaryKey;autoIncrement" json:"user_id"`
	Username  string    `gorm:"size:50;unique;not null" json:"username"`
	Password  string    `gorm:"size:255;not null" json:"password"`
	Email     string    `gorm:"size:100;unique;not null" json:"email"`
	CreatedAt time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP" json:"created_at"`
}

type ResetToken struct {
	ResetTokenID uint      `gorm:"primaryKey;autoIncrement" json:"reset_token_id"`
	Token        string    `gorm:"size:255;not null" json:"token"`
	Email        string    `gorm:"size:100;unique;not null" json:"email"`
	CreatedAt    time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP" json:"created_at"`
}
