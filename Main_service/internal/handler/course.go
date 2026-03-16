package handler

import (
	"net/http"
	"strconv"

	"github.com/baigel/lms/main-service/internal/model"
	"github.com/baigel/lms/main-service/internal/service"
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
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if err := h.service.CreateCourse(&course); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, course)
}

// GetCourseByID godoc
// @Summary Get course by ID
// @Description Get course by ID
// @Tags courses
// @Produce json
// @Param id path int true "Course ID"
// @Success 200 {object} model.Course
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /courses/{id} [get]
func (h *CourseHandler) GetCourseByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	course, err := h.service.GetCourseByID(uint(id))
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}

	c.JSON(http.StatusOK, course)
}

// GetAllCourses godoc
// @Summary Get all courses
// @Description Get all courses
// @Tags courses
// @Produce json
// @Success 200 {array} model.Course
// @Failure 500 {object} map[string]interface{}
// @Router /courses [get]
func (h *CourseHandler) GetAllCourses(c *gin.Context) {
	courses, err := h.service.GetAllCourses()
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, courses)
}

// UpdateCourse godoc
// @Summary Update a course
// @Description Update a course by ID
// @Tags courses
// @Accept json
// @Produce json
// @Param id path int true "Course ID"
// @Param course body model.Course true "Course data"
// @Success 200 {object} model.Course
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /courses/{id} [put]
func (h *CourseHandler) UpdateCourse(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var course model.Course
	if err := c.ShouldBindJSON(&course); err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}
	course.ID = uint(id)

	if err := h.service.UpdateCourse(&course); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, course)
}

// DeleteCourse godoc
// @Summary Delete a course
// @Description Delete a course by ID
// @Tags courses
// @Produce json
// @Param id path int true "Course ID"
// @Success 204
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /courses/{id} [delete]
func (h *CourseHandler) DeleteCourse(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	if err := h.service.DeleteCourse(uint(id)); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}
