package handler

import (
	"net/http"
	"strconv"

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

// CreateChapter godoc
// @Summary Create a new chapter
// @Description Create a new chapter
// @Tags chapters
// @Accept json
// @Produce json
// @Param chapter body model.Chapter true "Chapter data"
// @Success 201 {object} model.Chapter
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /chapters [post]
func (h *ChapterHandler) CreateChapter(c *gin.Context) {
	var chapter model.Chapter
	if err := c.ShouldBindJSON(&chapter); err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if err := h.service.CreateChapter(&chapter); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, chapter)
}


func (h *ChapterHandler) GetChapterByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	chapter, err := h.service.GetChapterByID(uint(id))
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Chapter not found"})
		return
	}

	c.JSON(http.StatusOK, chapter)
}


func (h *ChapterHandler) GetChaptersByCourseID(c *gin.Context) {
	courseID, err := strconv.ParseUint(c.Param("course_id"), 10, 32)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Course ID format"})
		return
	}

	chapters, err := h.service.GetChaptersByCourseID(uint(courseID))
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, chapters)
}

func (h *ChapterHandler) UpdateChapter(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var chapter model.Chapter
	if err := c.ShouldBindJSON(&chapter); err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}
	chapter.ID = uint(id)

	if err := h.service.UpdateChapter(&chapter); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, chapter)
}


func (h *ChapterHandler) DeleteChapter(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	if err := h.service.DeleteChapter(uint(id)); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}
