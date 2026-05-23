package model

import "time"

type Attachment struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"type:varchar(255);not null" json:"name"`
	URL       string    `gorm:"type:varchar(255);not null" json:"url"`
	LessonID  uint      `gorm:"not null;index" json:"lesson_id"`
	CreatedAt time.Time `json:"created_at"`
}
