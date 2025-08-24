package entities

import "time"

type TelegramUser struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	UserID     uint      `json:"user_id" gorm:"not null"`
	TelegramID int64     `json:"telegram_id" gorm:"uniqueIndex"`
	Username   string    `json:"username" gorm:"type:varchar(255)"`
	FirstName  string    `json:"first_name" gorm:"type:varchar(255)"`
	LastName   string    `json:"last_name" gorm:"type:varchar(255)"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// tablename
func (TelegramUser) TableName() string {
	return "telegram_users"
}
