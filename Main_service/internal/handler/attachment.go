package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/baigel/lms/main-service/internal/middleware"
	"github.com/baigel/lms/main-service/internal/service"
	"github.com/baigel/lms/main-service/pkg/apperror"
	"github.com/gin-gonic/gin"
)

type AttachmentHandler struct {
	service service.AttachmentService
}

func NewAttachmentHandler(service service.AttachmentService) *AttachmentHandler {
	return &AttachmentHandler{service: service}
}

func (h *AttachmentHandler) Upload(c *gin.Context) {
	lessonID, err := strconv.ParseUint(c.PostForm("lesson_id"), 10, 32)
	if err != nil {
		_ = c.Error(apperror.BadRequest("invalid lesson_id format", err))
		return
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		_ = c.Error(apperror.BadRequest("file is required", err))
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		_ = c.Error(apperror.BadRequest("failed to open uploaded file", err))
		return
	}
	defer file.Close()

	attachment, err := h.service.Upload(c.Request.Context(), service.AttachmentUpload{
		LessonID:    uint(lessonID),
		Name:        fileHeader.Filename,
		Reader:      file,
		Size:        fileHeader.Size,
		ContentType: fileHeader.Header.Get("Content-Type"),
	})
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, attachment)
}

func (h *AttachmentHandler) Download(c *gin.Context) {
	id, err := strconv.ParseUint(c.Query("id"), 10, 32)
	if err != nil {
		_ = c.Error(apperror.BadRequest("invalid id format", err))
		return
	}

	result, err := h.service.Download(c.Request.Context(), uint(id), middleware.RolesFromContext(c))
	if err != nil {
		_ = c.Error(err)
		return
	}
	defer result.File.Reader.Close()

	contentType := result.File.ContentType
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%q", result.ContentName))
	c.DataFromReader(http.StatusOK, result.File.Size, contentType, result.File.Reader, nil)
}
