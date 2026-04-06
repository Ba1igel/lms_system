package service

import (
	"errors"
	"testing"

	"github.com/baigel/lms/main-service/internal/model"
	"github.com/baigel/lms/main-service/internal/repository/mocks"
	"github.com/baigel/lms/main-service/pkg/apperror"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestLessonService_CreateLesson(t *testing.T) {
	mockRepo := mocks.NewLessonRepository(t)
	svc := NewLessonService(mockRepo)

	lesson := &model.Lesson{Name: "Test Lesson", ChapterID: 1}
	mockRepo.On("Create", lesson).Return(nil)

	err := svc.CreateLesson(lesson)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestLessonService_CreateLesson_RepoError(t *testing.T) {
	mockRepo := mocks.NewLessonRepository(t)
	svc := NewLessonService(mockRepo)

	lesson := &model.Lesson{Name: "Test Lesson", ChapterID: 1}
	mockRepo.On("Create", lesson).Return(errors.New("db error"))

	err := svc.CreateLesson(lesson)

	assert.Error(t, err)
	var appErr *apperror.AppError
	assert.True(t, errors.As(err, &appErr))
	assert.Equal(t, 500, appErr.Code)
	mockRepo.AssertExpectations(t)
}

func TestLessonService_GetLessonByID(t *testing.T) {
	mockRepo := mocks.NewLessonRepository(t)
	svc := NewLessonService(mockRepo)

	expectedLesson := &model.Lesson{ID: 1, Name: "Test Lesson"}
	mockRepo.On("GetByID", uint(1)).Return(expectedLesson, nil)

	lesson, err := svc.GetLessonByID(1)
	assert.NoError(t, err)
	assert.Equal(t, expectedLesson, lesson)
	mockRepo.AssertExpectations(t)
}

func TestLessonService_GetLessonByID_NotFound(t *testing.T) {
	mockRepo := mocks.NewLessonRepository(t)
	svc := NewLessonService(mockRepo)

	mockRepo.On("GetByID", uint(99)).Return(nil, gorm.ErrRecordNotFound)

	lesson, err := svc.GetLessonByID(99)

	assert.Nil(t, lesson)
	assert.Error(t, err)
	var appErr *apperror.AppError
	assert.True(t, errors.As(err, &appErr))
	assert.Equal(t, 404, appErr.Code)
	mockRepo.AssertExpectations(t)
}

func TestLessonService_GetLessonByID_RepoError(t *testing.T) {
	mockRepo := mocks.NewLessonRepository(t)
	svc := NewLessonService(mockRepo)

	mockRepo.On("GetByID", uint(1)).Return(nil, errors.New("db error"))

	lesson, err := svc.GetLessonByID(1)

	assert.Nil(t, lesson)
	assert.Error(t, err)
	var appErr *apperror.AppError
	assert.True(t, errors.As(err, &appErr))
	assert.Equal(t, 500, appErr.Code)
	mockRepo.AssertExpectations(t)
}

func TestLessonService_GetLessonsByChapterID(t *testing.T) {
	mockRepo := mocks.NewLessonRepository(t)
	svc := NewLessonService(mockRepo)

	expectedLessons := []model.Lesson{{ID: 1, Name: "Test Lesson", ChapterID: 1}}
	mockRepo.On("GetAllByChapterID", uint(1)).Return(expectedLessons, nil)

	lessons, err := svc.GetLessonsByChapterID(1)
	assert.NoError(t, err)
	assert.Equal(t, expectedLessons, lessons)
	mockRepo.AssertExpectations(t)
}

func TestLessonService_GetLessonsByChapterID_RepoError(t *testing.T) {
	mockRepo := mocks.NewLessonRepository(t)
	svc := NewLessonService(mockRepo)

	mockRepo.On("GetAllByChapterID", uint(1)).Return(nil, errors.New("db error"))

	lessons, err := svc.GetLessonsByChapterID(1)

	assert.Nil(t, lessons)
	assert.Error(t, err)
	var appErr *apperror.AppError
	assert.True(t, errors.As(err, &appErr))
	assert.Equal(t, 500, appErr.Code)
	mockRepo.AssertExpectations(t)
}

func TestLessonService_UpdateLesson(t *testing.T) {
	mockRepo := mocks.NewLessonRepository(t)
	svc := NewLessonService(mockRepo)

	lesson := &model.Lesson{ID: 1, Name: "Updated Lesson"}
	fresh := &model.Lesson{ID: 1, Name: "Updated Lesson", ChapterID: 2}
	mockRepo.On("Update", lesson).Return(nil)
	mockRepo.On("GetByID", uint(1)).Return(fresh, nil)

	result, err := svc.UpdateLesson(lesson)
	assert.NoError(t, err)
	assert.Equal(t, fresh, result)
	mockRepo.AssertExpectations(t)
}

func TestLessonService_UpdateLesson_NotFound(t *testing.T) {
	mockRepo := mocks.NewLessonRepository(t)
	svc := NewLessonService(mockRepo)

	lesson := &model.Lesson{ID: 99, Name: "Ghost Lesson"}
	mockRepo.On("Update", lesson).Return(gorm.ErrRecordNotFound)

	result, err := svc.UpdateLesson(lesson)

	assert.Nil(t, result)
	assert.Error(t, err)
	var appErr *apperror.AppError
	assert.True(t, errors.As(err, &appErr))
	assert.Equal(t, 404, appErr.Code)
	mockRepo.AssertExpectations(t)
}

func TestLessonService_UpdateLesson_RepoError(t *testing.T) {
	mockRepo := mocks.NewLessonRepository(t)
	svc := NewLessonService(mockRepo)

	lesson := &model.Lesson{ID: 1, Name: "Updated Lesson"}
	mockRepo.On("Update", lesson).Return(errors.New("db error"))

	result, err := svc.UpdateLesson(lesson)

	assert.Nil(t, result)
	assert.Error(t, err)
	var appErr *apperror.AppError
	assert.True(t, errors.As(err, &appErr))
	assert.Equal(t, 500, appErr.Code)
	mockRepo.AssertExpectations(t)
}

func TestLessonService_DeleteLesson(t *testing.T) {
	mockRepo := mocks.NewLessonRepository(t)
	svc := NewLessonService(mockRepo)

	mockRepo.On("Delete", uint(1)).Return(nil)

	err := svc.DeleteLesson(1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestLessonService_DeleteLesson_NotFound(t *testing.T) {
	mockRepo := mocks.NewLessonRepository(t)
	svc := NewLessonService(mockRepo)

	mockRepo.On("Delete", uint(99)).Return(gorm.ErrRecordNotFound)

	err := svc.DeleteLesson(99)

	assert.Error(t, err)
	var appErr *apperror.AppError
	assert.True(t, errors.As(err, &appErr))
	assert.Equal(t, 404, appErr.Code)
	mockRepo.AssertExpectations(t)
}

func TestLessonService_DeleteLesson_RepoError(t *testing.T) {
	mockRepo := mocks.NewLessonRepository(t)
	svc := NewLessonService(mockRepo)

	mockRepo.On("Delete", uint(1)).Return(errors.New("db error"))

	err := svc.DeleteLesson(1)

	assert.Error(t, err)
	var appErr *apperror.AppError
	assert.True(t, errors.As(err, &appErr))
	assert.Equal(t, 500, appErr.Code)
	mockRepo.AssertExpectations(t)
}
