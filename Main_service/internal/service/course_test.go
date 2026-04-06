package service

import (
	"testing"

	"github.com/baigel/lms/main-service/internal/model"
	"github.com/baigel/lms/main-service/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
)

func TestCourseService_CreateCourse(t *testing.T) {
	mockRepo := mocks.NewCourseRepository(t)
	svc := NewCourseService(mockRepo)

	course := &model.Course{Name: "Test Course"}
	mockRepo.On("Create", course).Return(nil)

	err := svc.CreateCourse(course)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestCourseService_GetCourseByID(t *testing.T) {
	mockRepo := mocks.NewCourseRepository(t)
	svc := NewCourseService(mockRepo)

	expectedCourse := &model.Course{ID: 1, Name: "Test Course"}
	mockRepo.On("GetByID", uint(1)).Return(expectedCourse, nil)

	course, err := svc.GetCourseByID(1)
	assert.NoError(t, err)
	assert.Equal(t, expectedCourse, course)
	mockRepo.AssertExpectations(t)
}

func TestCourseService_GetAllCourses(t *testing.T) {
	mockRepo := mocks.NewCourseRepository(t)
	svc := NewCourseService(mockRepo)

	expectedCourses := []model.Course{{ID: 1, Name: "Test Course"}}
	mockRepo.On("GetAll").Return(expectedCourses, nil)

	courses, err := svc.GetAllCourses()
	assert.NoError(t, err)
	assert.Equal(t, expectedCourses, courses)
	mockRepo.AssertExpectations(t)
}

func TestCourseService_UpdateCourse(t *testing.T) {
	mockRepo := mocks.NewCourseRepository(t)
	svc := NewCourseService(mockRepo)

	course := &model.Course{ID: 1, Name: "Updated Course"}
	mockRepo.On("Update", course).Return(nil)

	err := svc.UpdateCourse(course)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestCourseService_DeleteCourse(t *testing.T) {
	mockRepo := mocks.NewCourseRepository(t)
	svc := NewCourseService(mockRepo)

	mockRepo.On("Delete", uint(1)).Return(nil)

	err := svc.DeleteCourse(1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
