package handler

import (
	"net/http"
	"strconv"

	"github.com/baigel/lms/main-service/internal/dto"
	"github.com/baigel/lms/main-service/internal/model"
	"github.com/baigel/lms/main-service/internal/service"
	"github.com/gin-gonic/gin"
)

type LessonHandler struct {
	service service.LessonService
}

func NewLessonHandler(service service.LessonService) *LessonHandler {
	return &LessonHandler{service: service}
}

func lessonToResponse(l *model.Lesson) dto.LessonResponse {
	return dto.LessonResponse{
		ID:          l.ID,
		Name:        l.Name,
		Description: l.Description,
		Content:     l.Content,
		Order:       l.Order,
		ChapterID:   l.ChapterID,
		CreatedAt:   l.CreatedAt,
		UpdatedAt:   l.UpdatedAt,
	}
}

// CreateLesson godoc
// @Summary Create a new lesson
// @Description Create a new lesson within a chapter
// @Tags lessons
// @Accept json
// @Produce json
// @Param lesson body dto.CreateLessonRequest true "Lesson data"
// @Success 201 {object} dto.LessonResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /lessons [post]
func (h *LessonHandler) CreateLesson(c *gin.Context) {
	var req dto.CreateLessonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	lesson := &model.Lesson{
		Name:        req.Name,
		Description: req.Description,
		Content:     req.Content,
		Order:       req.Order,
		ChapterID:   req.ChapterID,
	}

	if err := h.service.CreateLesson(lesson); err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create lesson"})
		return
	}

	c.JSON(http.StatusCreated, lessonToResponse(lesson))
}

// GetLessonByID godoc
// @Summary Get lesson by ID
// @Description Get a single lesson by its ID
// @Tags lessons
// @Produce json
// @Param id path int true "Lesson ID"
// @Success 200 {object} dto.LessonResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /lessons/{id} [get]
func (h *LessonHandler) GetLessonByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	lesson, err := h.service.GetLessonByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Lesson not found"})
		return
	}

	c.JSON(http.StatusOK, lessonToResponse(lesson))
}

// GetLessonsByChapterID godoc
// @Summary Get lessons by chapter ID
// @Description Get all lessons belonging to a chapter
// @Tags lessons
// @Produce json
// @Param chapter_id path int true "Chapter ID"
// @Success 200 {array} dto.LessonResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /chapters/{chapter_id}/lessons [get]
func (h *LessonHandler) GetLessonsByChapterID(c *gin.Context) {
	chapterID, err := strconv.ParseUint(c.Param("chapter_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Chapter ID format"})
		return
	}

	lessons, err := h.service.GetLessonsByChapterID(uint(chapterID))
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get lessons"})
		return
	}

	resp := make([]dto.LessonResponse, len(lessons))
	for i := range lessons {
		resp[i] = lessonToResponse(&lessons[i])
	}

	c.JSON(http.StatusOK, resp)
}

// UpdateLesson godoc
// @Summary Update a lesson
// @Description Update an existing lesson by ID
// @Tags lessons
// @Accept json
// @Produce json
// @Param id path int true "Lesson ID"
// @Param lesson body dto.UpdateLessonRequest true "Lesson data"
// @Success 200 {object} dto.LessonResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /lessons/{id} [put]
func (h *LessonHandler) UpdateLesson(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var req dto.UpdateLessonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	lesson := &model.Lesson{
		ID:          uint(id),
		Name:        req.Name,
		Description: req.Description,
		Content:     req.Content,
		Order:       req.Order,
	}

	if err := h.service.UpdateLesson(lesson); err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update lesson"})
		return
	}

	c.JSON(http.StatusOK, lessonToResponse(lesson))
}

// DeleteLesson godoc
// @Summary Delete a lesson
// @Description Delete a lesson by ID
// @Tags lessons
// @Param id path int true "Lesson ID"
// @Success 204
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /lessons/{id} [delete]
func (h *LessonHandler) DeleteLesson(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	if err := h.service.DeleteLesson(uint(id)); err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete lesson"})
		return
	}

	c.Status(http.StatusNoContent)
}
