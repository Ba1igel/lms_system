package handler

import (
	"net/http"
	"strconv"

	"github.com/baigel/lms/main-service/internal/dto"
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

func courseToResponse(c *model.Course) dto.CourseResponse {
	return dto.CourseResponse{
		ID:          c.ID,
		Name:        c.Name,
		Description: c.Description,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}
}

// CreateCourse godoc
// @Summary Create a new course
// @Description Create a new course
// @Tags courses
// @Accept json
// @Produce json
// @Param course body dto.CreateCourseRequest true "Course data"
// @Success 201 {object} dto.CourseResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /courses [post]
func (h *CourseHandler) CreateCourse(c *gin.Context) {
	var req dto.CreateCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	course := &model.Course{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := h.service.CreateCourse(course); err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create course"})
		return
	}

	c.JSON(http.StatusCreated, courseToResponse(course))
}

// GetCourseByID godoc
// @Summary Get course by ID
// @Description Get a single course by its ID
// @Tags courses
// @Produce json
// @Param id path int true "Course ID"
// @Success 200 {object} dto.CourseResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /courses/{id} [get]
func (h *CourseHandler) GetCourseByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	course, err := h.service.GetCourseByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}

	c.JSON(http.StatusOK, courseToResponse(course))
}

// GetAllCourses godoc
// @Summary Get all courses
// @Description Get a list of all courses
// @Tags courses
// @Produce json
// @Success 200 {array} dto.CourseResponse
// @Failure 500 {object} map[string]interface{}
// @Router /courses [get]
func (h *CourseHandler) GetAllCourses(c *gin.Context) {
	courses, err := h.service.GetAllCourses()
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get courses"})
		return
	}

	resp := make([]dto.CourseResponse, len(courses))
	for i := range courses {
		resp[i] = courseToResponse(&courses[i])
	}

	c.JSON(http.StatusOK, resp)
}

// UpdateCourse godoc
// @Summary Update a course
// @Description Update an existing course by ID
// @Tags courses
// @Accept json
// @Produce json
// @Param id path int true "Course ID"
// @Param course body dto.UpdateCourseRequest true "Course data"
// @Success 200 {object} dto.CourseResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /courses/{id} [put]
func (h *CourseHandler) UpdateCourse(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var req dto.UpdateCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	course := &model.Course{
		ID:          uint(id),
		Name:        req.Name,
		Description: req.Description,
	}

	if err := h.service.UpdateCourse(course); err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update course"})
		return
	}

	c.JSON(http.StatusOK, courseToResponse(course))
}

// DeleteCourse godoc
// @Summary Delete a course
// @Description Delete a course by ID
// @Tags courses
// @Param id path int true "Course ID"
// @Success 204
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /courses/{id} [delete]
func (h *CourseHandler) DeleteCourse(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	if err := h.service.DeleteCourse(uint(id)); err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete course"})
		return
	}

	c.Status(http.StatusNoContent)
}
