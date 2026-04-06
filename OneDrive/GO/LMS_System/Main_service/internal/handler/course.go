package handler

import (
	"net/http"
	"strconv"

	"github.com/baigel/lms/main-service/internal/model"
	"github.com/baigel/lms/main-service/internal/service"
	"github.com/baigel/lms/main-service/pkg/apperror"
	"github.com/baigel/lms/main-service/pkg/pagination"
	"github.com/gin-gonic/gin"
)

type CourseHandler struct {
	service service.CourseService
}

func NewCourseHandler(service service.CourseService) *CourseHandler {
	return &CourseHandler{service: service}
}

// CreateCourse godoc
// @Summary Create a new course
// @Description Create a new course
// @Tags courses
// @Accept json
// @Produce json
// @Param course body model.Course true "Course data"
// @Success 201 {object} model.Course
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /courses [post]
func (h *CourseHandler) CreateCourse(c *gin.Context) {
	var course model.Course
	if err := c.ShouldBindJSON(&course); err != nil {
		_ = c.Error(apperror.BadRequest("invalid request payload", err))
		return
	}

	if err := h.service.CreateCourse(&course); err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, course)
}

func (h *CourseHandler) GetCourseByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		_ = c.Error(apperror.BadRequest("invalid id format", err))
		return
	}

	course, err := h.service.GetCourseByID(uint(id))
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, course)
}

func (h *CourseHandler) GetAllCourses(c *gin.Context) {
	p := pagination.FromContext(c)

	courses, total, err := h.service.GetAllCourses(p)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, pagination.NewResponse(courses, total, p))
}

func (h *CourseHandler) UpdateCourse(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		_ = c.Error(apperror.BadRequest("invalid id format", err))
		return
	}

	var course model.Course
	if err := c.ShouldBindJSON(&course); err != nil {
		_ = c.Error(apperror.BadRequest("invalid request payload", err))
		return
	}
	course.ID = uint(id)

	updated, err := h.service.UpdateCourse(&course)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, updated)
}

func (h *CourseHandler) DeleteCourse(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		_ = c.Error(apperror.BadRequest("invalid id format", err))
		return
	}

	if err := h.service.DeleteCourse(uint(id)); err != nil {
		_ = c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}
