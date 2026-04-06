package service

import (
	"errors"
	"testing"

	"github.com/baigel/lms/main-service/internal/model"
	"github.com/baigel/lms/main-service/internal/repository/mocks"
	"github.com/baigel/lms/main-service/pkg/apperror"
	"github.com/baigel/lms/main-service/pkg/pagination"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
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

func TestCourseService_CreateCourse_RepoError(t *testing.T) {
	mockRepo := mocks.NewCourseRepository(t)
	svc := NewCourseService(mockRepo)

	course := &model.Course{Name: "Test Course"}
	mockRepo.On("Create", course).Return(errors.New("db error"))

	err := svc.CreateCourse(course)

	assert.Error(t, err)
	var appErr *apperror.AppError
	assert.True(t, errors.As(err, &appErr))
	assert.Equal(t, 500, appErr.Code)
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

func TestCourseService_GetCourseByID_NotFound(t *testing.T) {
	mockRepo := mocks.NewCourseRepository(t)
	svc := NewCourseService(mockRepo)

	mockRepo.On("GetByID", uint(99)).Return(nil, gorm.ErrRecordNotFound)

	course, err := svc.GetCourseByID(99)

	assert.Nil(t, course)
	assert.Error(t, err)
	var appErr *apperror.AppError
	assert.True(t, errors.As(err, &appErr))
	assert.Equal(t, 404, appErr.Code)
	mockRepo.AssertExpectations(t)
}

func TestCourseService_GetCourseByID_RepoError(t *testing.T) {
	mockRepo := mocks.NewCourseRepository(t)
	svc := NewCourseService(mockRepo)

	mockRepo.On("GetByID", uint(1)).Return(nil, errors.New("db error"))

	course, err := svc.GetCourseByID(1)

	assert.Nil(t, course)
	assert.Error(t, err)
	var appErr *apperror.AppError
	assert.True(t, errors.As(err, &appErr))
	assert.Equal(t, 500, appErr.Code)
	mockRepo.AssertExpectations(t)
}

func TestCourseService_GetAllCourses(t *testing.T) {
	mockRepo := mocks.NewCourseRepository(t)
	svc := NewCourseService(mockRepo)

	p := pagination.Params{Page: 1, Limit: 20}
	expectedCourses := []model.Course{{ID: 1, Name: "Test Course"}}
	mockRepo.On("GetAll", p).Return(expectedCourses, int64(1), nil)

	courses, total, err := svc.GetAllCourses(p)
	assert.NoError(t, err)
	assert.Equal(t, expectedCourses, courses)
	assert.Equal(t, int64(1), total)
	mockRepo.AssertExpectations(t)
}

func TestCourseService_GetAllCourses_RepoError(t *testing.T) {
	mockRepo := mocks.NewCourseRepository(t)
	svc := NewCourseService(mockRepo)

	p := pagination.Params{Page: 1, Limit: 20}
	mockRepo.On("GetAll", p).Return(nil, int64(0), errors.New("db error"))

	courses, total, err := svc.GetAllCourses(p)

	assert.Nil(t, courses)
	assert.Equal(t, int64(0), total)
	assert.Error(t, err)
	var appErr *apperror.AppError
	assert.True(t, errors.As(err, &appErr))
	assert.Equal(t, 500, appErr.Code)
	mockRepo.AssertExpectations(t)
}

func TestCourseService_UpdateCourse(t *testing.T) {
	mockRepo := mocks.NewCourseRepository(t)
	svc := NewCourseService(mockRepo)

	course := &model.Course{ID: 1, Name: "Updated Course"}
	fresh := &model.Course{ID: 1, Name: "Updated Course", Description: "from db"}
	mockRepo.On("Update", course).Return(nil)
	mockRepo.On("GetByID", uint(1)).Return(fresh, nil)

	result, err := svc.UpdateCourse(course)
	assert.NoError(t, err)
	assert.Equal(t, fresh, result)
	mockRepo.AssertExpectations(t)
}

func TestCourseService_UpdateCourse_NotFound(t *testing.T) {
	mockRepo := mocks.NewCourseRepository(t)
	svc := NewCourseService(mockRepo)

	course := &model.Course{ID: 99, Name: "Ghost Course"}
	mockRepo.On("Update", course).Return(gorm.ErrRecordNotFound)

	result, err := svc.UpdateCourse(course)

	assert.Nil(t, result)
	assert.Error(t, err)
	var appErr *apperror.AppError
	assert.True(t, errors.As(err, &appErr))
	assert.Equal(t, 404, appErr.Code)
	mockRepo.AssertExpectations(t)
}

func TestCourseService_UpdateCourse_RepoError(t *testing.T) {
	mockRepo := mocks.NewCourseRepository(t)
	svc := NewCourseService(mockRepo)

	course := &model.Course{ID: 1, Name: "Updated Course"}
	mockRepo.On("Update", course).Return(errors.New("db error"))

	result, err := svc.UpdateCourse(course)

	assert.Nil(t, result)
	assert.Error(t, err)
	var appErr *apperror.AppError
	assert.True(t, errors.As(err, &appErr))
	assert.Equal(t, 500, appErr.Code)
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

func TestCourseService_DeleteCourse_NotFound(t *testing.T) {
	mockRepo := mocks.NewCourseRepository(t)
	svc := NewCourseService(mockRepo)

	mockRepo.On("Delete", uint(99)).Return(gorm.ErrRecordNotFound)

	err := svc.DeleteCourse(99)

	assert.Error(t, err)
	var appErr *apperror.AppError
	assert.True(t, errors.As(err, &appErr))
	assert.Equal(t, 404, appErr.Code)
	mockRepo.AssertExpectations(t)
}

func TestCourseService_DeleteCourse_RepoError(t *testing.T) {
	mockRepo := mocks.NewCourseRepository(t)
	svc := NewCourseService(mockRepo)

	mockRepo.On("Delete", uint(1)).Return(errors.New("db error"))

	err := svc.DeleteCourse(1)

	assert.Error(t, err)
	var appErr *apperror.AppError
	assert.True(t, errors.As(err, &appErr))
	assert.Equal(t, 500, appErr.Code)
	mockRepo.AssertExpectations(t)
}
