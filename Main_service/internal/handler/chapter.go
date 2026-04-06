package handler

import (
	"net/http"
	"strconv"

	"github.com/baigel/lms/main-service/internal/model"
	"github.com/baigel/lms/main-service/internal/service"
	"github.com/baigel/lms/main-service/pkg/apperror"
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
		_ = c.Error(apperror.BadRequest("invalid request payload", err))
		return
	}

	if err := h.service.CreateChapter(&chapter); err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, chapter)
}

func (h *ChapterHandler) GetChapterByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		_ = c.Error(apperror.BadRequest("invalid id format", err))
		return
	}

	chapter, err := h.service.GetChapterByID(uint(id))
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, chapter)
}

func (h *ChapterHandler) GetChaptersByCourseID(c *gin.Context) {
	courseID, err := strconv.ParseUint(c.Param("course_id"), 10, 32)
	if err != nil {
		_ = c.Error(apperror.BadRequest("invalid course id format", err))
		return
	}

	chapters, err := h.service.GetChaptersByCourseID(uint(courseID))
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, chapters)
}

func (h *ChapterHandler) UpdateChapter(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		_ = c.Error(apperror.BadRequest("invalid id format", err))
		return
	}

	var chapter model.Chapter
	if err := c.ShouldBindJSON(&chapter); err != nil {
		_ = c.Error(apperror.BadRequest("invalid request payload", err))
		return
	}
	chapter.ID = uint(id)

	updated, err := h.service.UpdateChapter(&chapter)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, updated)
}

func (h *ChapterHandler) DeleteChapter(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		_ = c.Error(apperror.BadRequest("invalid id format", err))
		return
	}

	if err := h.service.DeleteChapter(uint(id)); err != nil {
		_ = c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}
