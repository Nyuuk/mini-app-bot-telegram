package entities

import "time"

type Overtime struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	TelegramUserID  uint      `json:"telegram_user_id" gorm:"default:null"`
	Date            time.Time `json:"date" gorm:"type:date;not null"`                      // Format: YYYY-MM-DD
	TimeStart       time.Time `json:"time_start" gorm:"type:time;defualt:null"`            // Format: HH:MM
	TimeStop        time.Time `json:"time_stop" gorm:"type:time;defualt:null"`             // Format: HH:MM
	BreakDuration   float64   `json:"break_duration" gorm:"type:decimal(4,2);default:0.0"` // Durasi istirahat dalam jam (1.5 = 1 jam 30 menit)
	Duration        float64   `json:"duration" gorm:"type:decimal(4,2);default:0.0"`       // Durasi total lembur dikurangi break (manual input)
	Description     string    `json:"description" gorm:"type:text;default:null"`
	Category        string    `json:"category" gorm:"type:varchar(255);default:Other"`
	CreatedByUserID uint      `json:"-" gorm:"not null"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	CreatedAt       time.Time `json:"created_at" gorm:"autoCreateTime"`

	// Relasi
	User         User         `json:"-" gorm:"foreignKey:CreatedByUserID"`
	TelegramUser TelegramUser `json:"-" gorm:"foreignKey:TelegramUserID"`
}

// tablename
func (Overtime) TableName() string {
	return "overtimes"
}
