package dto

import "time"

type CreateLessonRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Content     string `json:"content"`
	Order       int    `json:"order" binding:"required"`
	ChapterID   uint   `json:"chapter_id" binding:"required"`
}

type UpdateLessonRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Content     string `json:"content"`
	Order       int    `json:"order" binding:"required"`
}

type LessonResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Content     string    `json:"content"`
	Order       int       `json:"order"`
	ChapterID   uint      `json:"chapter_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
