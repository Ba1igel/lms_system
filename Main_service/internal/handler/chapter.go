package handler

import (
	"net/http"
	"strconv"

	"github.com/baigel/lms/main-service/internal/dto"
	"github.com/baigel/lms/main-service/internal/model"
	"github.com/baigel/lms/main-service/internal/service"
	"github.com/gin-gonic/gin"
)

type ChapterHandler struct {
	service service.ChapterService
}

func NewChapterHandler(service service.ChapterService) *ChapterHandler {
	return &ChapterHandler{service: service}
}

func chapterToResponse(ch *model.Chapter) dto.ChapterResponse {
	return dto.ChapterResponse{
		ID:          ch.ID,
		Name:        ch.Name,
		Description: ch.Description,
		Order:       ch.Order,
		CourseID:    ch.CourseID,
		CreatedAt:   ch.CreatedAt,
		UpdatedAt:   ch.UpdatedAt,
	}
}

// CreateChapter godoc
// @Summary Create a new chapter
// @Description Create a new chapter within a course
// @Tags chapters
// @Accept json
// @Produce json
// @Param chapter body dto.CreateChapterRequest true "Chapter data"
// @Success 201 {object} dto.ChapterResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /chapters [post]
func (h *ChapterHandler) CreateChapter(c *gin.Context) {
	var req dto.CreateChapterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	chapter := &model.Chapter{
		Name:        req.Name,
		Description: req.Description,
		Order:       req.Order,
		CourseID:    req.CourseID,
	}

	if err := h.service.CreateChapter(chapter); err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create chapter"})
		return
	}

	c.JSON(http.StatusCreated, chapterToResponse(chapter))
}

// GetChapterByID godoc
// @Summary Get chapter by ID
// @Description Get a single chapter by its ID
// @Tags chapters
// @Produce json
// @Param id path int true "Chapter ID"
// @Success 200 {object} dto.ChapterResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /chapters/{id} [get]
func (h *ChapterHandler) GetChapterByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	chapter, err := h.service.GetChapterByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Chapter not found"})
		return
	}

	c.JSON(http.StatusOK, chapterToResponse(chapter))
}

// GetChaptersByCourseID godoc
// @Summary Get chapters by course ID
// @Description Get all chapters belonging to a course
// @Tags chapters
// @Produce json
// @Param course_id path int true "Course ID"
// @Success 200 {array} dto.ChapterResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /courses/{course_id}/chapters [get]
func (h *ChapterHandler) GetChaptersByCourseID(c *gin.Context) {
	courseID, err := strconv.ParseUint(c.Param("course_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Course ID format"})
		return
	}

	chapters, err := h.service.GetChaptersByCourseID(uint(courseID))
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get chapters"})
		return
	}

	resp := make([]dto.ChapterResponse, len(chapters))
	for i := range chapters {
		resp[i] = chapterToResponse(&chapters[i])
	}

	c.JSON(http.StatusOK, resp)
}

// UpdateChapter godoc
// @Summary Update a chapter
// @Description Update an existing chapter by ID
// @Tags chapters
// @Accept json
// @Produce json
// @Param id path int true "Chapter ID"
// @Param chapter body dto.UpdateChapterRequest true "Chapter data"
// @Success 200 {object} dto.ChapterResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /chapters/{id} [put]
func (h *ChapterHandler) UpdateChapter(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var req dto.UpdateChapterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	chapter := &model.Chapter{
		ID:          uint(id),
		Name:        req.Name,
		Description: req.Description,
		Order:       req.Order,
	}

	if err := h.service.UpdateChapter(chapter); err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update chapter"})
		return
	}

	c.JSON(http.StatusOK, chapterToResponse(chapter))
}

// DeleteChapter godoc
// @Summary Delete a chapter
// @Description Delete a chapter by ID
// @Tags chapters
// @Param id path int true "Chapter ID"
// @Success 204
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /chapters/{id} [delete]
func (h *ChapterHandler) DeleteChapter(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	if err := h.service.DeleteChapter(uint(id)); err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete chapter"})
		return
	}

	c.Status(http.StatusNoContent)
}
