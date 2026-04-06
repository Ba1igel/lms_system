package handler

import (
	"net/http"
	"strconv"

	"github.com/baigel/lms/main-service/internal/model"
	"github.com/baigel/lms/main-service/internal/service"
	"github.com/baigel/lms/main-service/pkg/apperror"
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
		_ = c.Error(apperror.BadRequest("invalid request payload", err))
		return
	}

	if err := h.service.CreateLesson(&lesson); err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, lesson)
}

func (h *LessonHandler) GetLessonByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		_ = c.Error(apperror.BadRequest("invalid id format", err))
		return
	}

	lesson, err := h.service.GetLessonByID(uint(id))
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, lesson)
}

func (h *LessonHandler) GetLessonsByChapterID(c *gin.Context) {
	chapterID, err := strconv.ParseUint(c.Param("chapter_id"), 10, 32)
	if err != nil {
		_ = c.Error(apperror.BadRequest("invalid chapter id format", err))
		return
	}

	lessons, err := h.service.GetLessonsByChapterID(uint(chapterID))
	if err != nil {
		_ = c.Error(err)
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
		_ = c.Error(apperror.BadRequest("invalid id format", err))
		return
	}

	var lesson model.Lesson
	if err := c.ShouldBindJSON(&lesson); err != nil {
		_ = c.Error(apperror.BadRequest("invalid request payload", err))
		return
	}
	lesson.ID = uint(id)

	updated, err := h.service.UpdateLesson(&lesson)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, updated)
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
		_ = c.Error(apperror.BadRequest("invalid id format", err))
		return
	}

	if err := h.service.DeleteLesson(uint(id)); err != nil {
		_ = c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}
