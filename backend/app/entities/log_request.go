package entities

import (
	"time"
)

type LogRequest struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	UserID       uint      `json:"user_id" gorm:"not null"`
	Method       string    `json:"method" gorm:"type:varchar(255);not null"`
	Endpoint     string    `json:"endpoint" gorm:"type:varchar(255);not null"`
	RequestBody  string    `json:"request_body" gorm:"type:text;not null"`
	ResponseBody string    `json:"response_body" gorm:"type:text;not null"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
	// Relasi
	// User User `json:"-" gorm:"foreignKey:UserID"`
}

// tablename
func (LogRequest) TableName() string {
	return "log_requests"
}
