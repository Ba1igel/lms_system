package dto

import "time"

type CreateCourseRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type UpdateCourseRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type CourseResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
