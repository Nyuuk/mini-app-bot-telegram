package entities

import "time"

type User struct {
	ID           uint         `json:"id" gorm:"primaryKey"`
	Username     string       `json:"username" gorm:"uniqueIndex"`
	PasswordHash string       `json:"-" gorm:"not null"`
	Email        string       `json:"email" gorm:"uniqueIndex"`
	CreatedAt    time.Time    `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time    `json:"updated_at" gorm:"autoUpdateTime"`
	APIKeys      []APIKey     `json:"-"` // relasi one-to-many
	TelegramUser TelegramUser `json:"-"` // relasi one-to-one
}

// tablename
func (User) TableName() string {
	return "users"
}
