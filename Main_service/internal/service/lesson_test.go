package service

import (
	"testing"

	"github.com/baigel/lms/main-service/internal/model"
	"github.com/baigel/lms/main-service/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
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

func TestLessonService_UpdateLesson(t *testing.T) {
	mockRepo := mocks.NewLessonRepository(t)
	svc := NewLessonService(mockRepo)

	lesson := &model.Lesson{ID: 1, Name: "Updated Lesson"}
	mockRepo.On("Update", lesson).Return(nil)

	err := svc.UpdateLesson(lesson)
	assert.NoError(t, err)
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
