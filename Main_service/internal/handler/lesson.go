 package handler

import (
	"net/http"
	"strconv"

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

// CreateLesson godoc
// @Summary Create a new lesson
// @Description Create a new lesson
// @Tags lessons
// @Accept json
// @Produce json
// @Param lesson body model.Lesson true "Lesson data"
// @Success 201 {object} model.Lesson
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /lessons [post]
func (h *LessonHandler) CreateLesson(c *gin.Context) {
	var lesson model.Lesson
	if err := c.ShouldBindJSON(&lesson); err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if err := h.service.CreateLesson(&lesson); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, lesson)
}

// GetLessonByID godoc
// @Summary Get lesson by ID
// @Description Get lesson by ID
// @Tags lessons
// @Produce json
// @Param id path int true "Lesson ID"
// @Success 200 {object} model.Lesson
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /lessons/{id} [get]
func (h *LessonHandler) GetLessonByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	lesson, err := h.service.GetLessonByID(uint(id))
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Lesson not found"})
		return
	}

	c.JSON(http.StatusOK, lesson)
}

// GetLessonsByChapterID godoc
// @Summary Get lessons by Chapter ID
// @Description Get all lessons associated with a specific chapter ID
// @Tags lessons
// @Produce json
// @Param chapter_id path int true "Chapter ID"
// @Success 200 {array} model.Lesson
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /chapters/{chapter_id}/lessons [get]
func (h *LessonHandler) GetLessonsByChapterID(c *gin.Context) {
	chapterID, err := strconv.ParseUint(c.Param("chapter_id"), 10, 32)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Chapter ID format"})
		return
	}

	lessons, err := h.service.GetLessonsByChapterID(uint(chapterID))
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, lessons)
}

// UpdateLesson godoc
// @Summary Update a lesson
// @Description Update a lesson by ID
// @Tags lessons
// @Accept json
// @Produce json
// @Param id path int true "Lesson ID"
// @Param lesson body model.Lesson true "Lesson data"
// @Success 200 {object} model.Lesson
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /lessons/{id} [put]
func (h *LessonHandler) UpdateLesson(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var lesson model.Lesson
	if err := c.ShouldBindJSON(&lesson); err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}
	lesson.ID = uint(id)

	if err := h.service.UpdateLesson(&lesson); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, lesson)
}

// DeleteLesson godoc
// @Summary Delete a lesson
// @Description Delete a lesson by ID
// @Tags lessons
// @Produce json
// @Param id path int true "Lesson ID"
// @Success 204
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /lessons/{id} [delete]
func (h *LessonHandler) DeleteLesson(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	if err := h.service.DeleteLesson(uint(id)); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}
