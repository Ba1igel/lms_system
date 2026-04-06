package service

import (
	"errors"

	"github.com/baigel/lms/main-service/internal/model"
	"github.com/baigel/lms/main-service/internal/repository"
	"github.com/baigel/lms/main-service/pkg/apperror"
	"gorm.io/gorm"
)

type LessonService interface {
	CreateLesson(lesson *model.Lesson) error
	GetLessonByID(id uint) (*model.Lesson, error)
	GetLessonsByChapterID(chapterID uint) ([]model.Lesson, error)
	UpdateLesson(lesson *model.Lesson) (*model.Lesson, error)
	DeleteLesson(id uint) error
}

type lessonService struct {
	repo repository.LessonRepository
}

func NewLessonService(repo repository.LessonRepository) LessonService {
	return &lessonService{repo: repo}
}

func (s *lessonService) CreateLesson(lesson *model.Lesson) error {
	if err := s.repo.Create(lesson); err != nil {
		return apperror.Internal("failed to create lesson", err)
	}
	return nil
}

func (s *lessonService) GetLessonByID(id uint) (*model.Lesson, error) {
	lesson, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFound("lesson not found")
		}
		return nil, apperror.Internal("failed to get lesson", err)
	}
	return lesson, nil
}

func (s *lessonService) GetLessonsByChapterID(chapterID uint) ([]model.Lesson, error) {
	lessons, err := s.repo.GetAllByChapterID(chapterID)
	if err != nil {
		return nil, apperror.Internal("failed to get lessons", err)
	}
	return lessons, nil
}

func (s *lessonService) UpdateLesson(lesson *model.Lesson) (*model.Lesson, error) {
	if err := s.repo.Update(lesson); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFound("lesson not found")
		}
		return nil, apperror.Internal("failed to update lesson", err)
	}

	updated, err := s.repo.GetByID(lesson.ID)
	if err != nil {
		return nil, apperror.Internal("failed to fetch updated lesson", err)
	}
	return updated, nil
}

func (s *lessonService) DeleteLesson(id uint) error {
	if err := s.repo.Delete(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.NotFound("lesson not found")
		}
		return apperror.Internal("failed to delete lesson", err)
	}
	return nil
}
