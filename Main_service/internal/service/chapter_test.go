package service

import (
	"testing"

	"github.com/baigel/lms/main-service/internal/model"
	"github.com/baigel/lms/main-service/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
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

func TestChapterService_UpdateChapter(t *testing.T) {
	mockRepo := mocks.NewChapterRepository(t)
	svc := NewChapterService(mockRepo)

	chapter := &model.Chapter{ID: 1, Name: "Updated Chapter"}
	mockRepo.On("Update", chapter).Return(nil)

	err := svc.UpdateChapter(chapter)
	assert.NoError(t, err)
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
