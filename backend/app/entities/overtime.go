package entities

import (
	"encoding/json"
	"time"
)

type Overtime struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	TelegramUserID  uint      `json:"telegram_user_id" gorm:"default:null"`
	Date            time.Time `json:"date" gorm:"type:date;not null"`                      // Format: YYYY-MM-DD
	TimeStart       time.Time `json:"-" gorm:"type:time;defualt:null"`                     // Format: HH:MM (custom JSON)
	TimeStop        time.Time `json:"-" gorm:"type:time;defualt:null"`                     // Format: HH:MM (custom JSON)
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

// Custom JSON marshaling for time-only fields
func (o Overtime) MarshalJSON() ([]byte, error) {
	type Alias Overtime

	// Format time-only fields
	timeStartStr := ""
	timeStopStr := ""

	// Check if it's a time-only value (base date 1900-01-01)
	if o.TimeStart.Year() == 1900 {
		timeStartStr = o.TimeStart.Format("15:04:05")
	} else {
		timeStartStr = o.TimeStart.Format("2006-01-02T15:04:05Z07:00")
	}

	if o.TimeStop.Year() == 1900 {
		timeStopStr = o.TimeStop.Format("15:04:05")
	} else {
		timeStopStr = o.TimeStop.Format("2006-01-02T15:04:05Z07:00")
	}

	return json.Marshal(&struct {
		TimeStart string `json:"time_start"`
		TimeStop  string `json:"time_stop"`
		*Alias
	}{
		TimeStart: timeStartStr,
		TimeStop:  timeStopStr,
		Alias:     (*Alias)(&o),
	})
}

// tablename
func (Overtime) TableName() string {
	return "overtimes"
}
