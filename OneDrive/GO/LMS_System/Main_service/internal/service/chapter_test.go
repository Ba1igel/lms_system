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

func TestChapterService_CreateChapter(t *testing.T) {
	mockRepo := mocks.NewChapterRepository(t)
	svc := NewChapterService(mockRepo)

	chapter := &model.Chapter{Name: "Test Chapter", CourseID: 1}
	mockRepo.On("Create", chapter).Return(nil)

	err := svc.CreateChapter(chapter)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestChapterService_CreateChapter_RepoError(t *testing.T) {
	mockRepo := mocks.NewChapterRepository(t)
	svc := NewChapterService(mockRepo)

	chapter := &model.Chapter{Name: "Test Chapter", CourseID: 1}
	mockRepo.On("Create", chapter).Return(errors.New("db error"))

	err := svc.CreateChapter(chapter)

	assert.Error(t, err)
	var appErr *apperror.AppError
	assert.True(t, errors.As(err, &appErr))
	assert.Equal(t, 500, appErr.Code)
	mockRepo.AssertExpectations(t)
}

func TestChapterService_GetChapterByID(t *testing.T) {
	mockRepo := mocks.NewChapterRepository(t)
	svc := NewChapterService(mockRepo)

	expectedChapter := &model.Chapter{ID: 1, Name: "Test Chapter"}
	mockRepo.On("GetByID", uint(1)).Return(expectedChapter, nil)

	chapter, err := svc.GetChapterByID(1)
	assert.NoError(t, err)
	assert.Equal(t, expectedChapter, chapter)
	mockRepo.AssertExpectations(t)
}

func TestChapterService_GetChapterByID_NotFound(t *testing.T) {
	mockRepo := mocks.NewChapterRepository(t)
	svc := NewChapterService(mockRepo)

	mockRepo.On("GetByID", uint(99)).Return(nil, gorm.ErrRecordNotFound)

	chapter, err := svc.GetChapterByID(99)

	assert.Nil(t, chapter)
	assert.Error(t, err)
	var appErr *apperror.AppError
	assert.True(t, errors.As(err, &appErr))
	assert.Equal(t, 404, appErr.Code)
	mockRepo.AssertExpectations(t)
}

func TestChapterService_GetChapterByID_RepoError(t *testing.T) {
	mockRepo := mocks.NewChapterRepository(t)
	svc := NewChapterService(mockRepo)

	mockRepo.On("GetByID", uint(1)).Return(nil, errors.New("db error"))

	chapter, err := svc.GetChapterByID(1)

	assert.Nil(t, chapter)
	assert.Error(t, err)
	var appErr *apperror.AppError
	assert.True(t, errors.As(err, &appErr))
	assert.Equal(t, 500, appErr.Code)
	mockRepo.AssertExpectations(t)
}

func TestChapterService_GetChaptersByCourseID(t *testing.T) {
	mockRepo := mocks.NewChapterRepository(t)
	svc := NewChapterService(mockRepo)

	expectedChapters := []model.Chapter{{ID: 1, Name: "Test Chapter", CourseID: 1}}
	mockRepo.On("GetAllByCourseID", uint(1)).Return(expectedChapters, nil)

	chapters, err := svc.GetChaptersByCourseID(1)
	assert.NoError(t, err)
	assert.Equal(t, expectedChapters, chapters)
	mockRepo.AssertExpectations(t)
}

func TestChapterService_GetChaptersByCourseID_RepoError(t *testing.T) {
	mockRepo := mocks.NewChapterRepository(t)
	svc := NewChapterService(mockRepo)

	mockRepo.On("GetAllByCourseID", uint(1)).Return(nil, errors.New("db error"))

	chapters, err := svc.GetChaptersByCourseID(1)

	assert.Nil(t, chapters)
	assert.Error(t, err)
	var appErr *apperror.AppError
	assert.True(t, errors.As(err, &appErr))
	assert.Equal(t, 500, appErr.Code)
	mockRepo.AssertExpectations(t)
}

func TestChapterService_UpdateChapter(t *testing.T) {
	mockRepo := mocks.NewChapterRepository(t)
	svc := NewChapterService(mockRepo)

	chapter := &model.Chapter{ID: 1, Name: "Updated Chapter"}
	fresh := &model.Chapter{ID: 1, Name: "Updated Chapter", CourseID: 2}
	mockRepo.On("Update", chapter).Return(nil)
	mockRepo.On("GetByID", uint(1)).Return(fresh, nil)

	result, err := svc.UpdateChapter(chapter)
	assert.NoError(t, err)
	assert.Equal(t, fresh, result)
	mockRepo.AssertExpectations(t)
}

func TestChapterService_UpdateChapter_NotFound(t *testing.T) {
	mockRepo := mocks.NewChapterRepository(t)
	svc := NewChapterService(mockRepo)

	chapter := &model.Chapter{ID: 99, Name: "Ghost Chapter"}
	mockRepo.On("Update", chapter).Return(gorm.ErrRecordNotFound)

	result, err := svc.UpdateChapter(chapter)

	assert.Nil(t, result)
	assert.Error(t, err)
	var appErr *apperror.AppError
	assert.True(t, errors.As(err, &appErr))
	assert.Equal(t, 404, appErr.Code)
	mockRepo.AssertExpectations(t)
}

func TestChapterService_UpdateChapter_RepoError(t *testing.T) {
	mockRepo := mocks.NewChapterRepository(t)
	svc := NewChapterService(mockRepo)

	chapter := &model.Chapter{ID: 1, Name: "Updated Chapter"}
	mockRepo.On("Update", chapter).Return(errors.New("db error"))

	result, err := svc.UpdateChapter(chapter)

	assert.Nil(t, result)
	assert.Error(t, err)
	var appErr *apperror.AppError
	assert.True(t, errors.As(err, &appErr))
	assert.Equal(t, 500, appErr.Code)
	mockRepo.AssertExpectations(t)
}

func TestChapterService_DeleteChapter(t *testing.T) {
	mockRepo := mocks.NewChapterRepository(t)
	svc := NewChapterService(mockRepo)

	mockRepo.On("Delete", uint(1)).Return(nil)

	err := svc.DeleteChapter(1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestChapterService_DeleteChapter_NotFound(t *testing.T) {
	mockRepo := mocks.NewChapterRepository(t)
	svc := NewChapterService(mockRepo)

	mockRepo.On("Delete", uint(99)).Return(gorm.ErrRecordNotFound)

	err := svc.DeleteChapter(99)

	assert.Error(t, err)
	var appErr *apperror.AppError
	assert.True(t, errors.As(err, &appErr))
	assert.Equal(t, 404, appErr.Code)
	mockRepo.AssertExpectations(t)
}

func TestChapterService_DeleteChapter_RepoError(t *testing.T) {
	mockRepo := mocks.NewChapterRepository(t)
	svc := NewChapterService(mockRepo)

	mockRepo.On("Delete", uint(1)).Return(errors.New("db error"))

	err := svc.DeleteChapter(1)

	assert.Error(t, err)
	var appErr *apperror.AppError
	assert.True(t, errors.As(err, &appErr))
	assert.Equal(t, 500, appErr.Code)
	mockRepo.AssertExpectations(t)
}
