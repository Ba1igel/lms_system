package model

import "time"

type Chapter struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"type:varchar(255);not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Order       int       `gorm:"not null" json:"order"`
	CourseID    uint      `gorm:"not null" json:"course_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
