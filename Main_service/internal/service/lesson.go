package service

import (
	"github.com/baigel/lms/main-service/internal/model"
	"github.com/baigel/lms/main-service/internal/repository"
)

type LessonService interface {
	CreateLesson(lesson *model.Lesson) error
	GetLessonByID(id uint) (*model.Lesson, error)
	GetLessonsByChapterID(chapterID uint) ([]model.Lesson, error)
	UpdateLesson(lesson *model.Lesson) error
	DeleteLesson(id uint) error
}

type lessonService struct {
	repo repository.LessonRepository
}

func NewLessonService(repo repository.LessonRepository) LessonService {
	return &lessonService{repo: repo}
}

func (s *lessonService) CreateLesson(lesson *model.Lesson) error {
	return s.repo.Create(lesson)
}

func (s *lessonService) GetLessonByID(id uint) (*model.Lesson, error) {
	return s.repo.GetByID(id)
}

func (s *lessonService) GetLessonsByChapterID(chapterID uint) ([]model.Lesson, error) {
	return s.repo.GetAllByChapterID(chapterID)
}

func (s *lessonService) UpdateLesson(lesson *model.Lesson) error {
	return s.repo.Update(lesson)
}

func (s *lessonService) DeleteLesson(id uint) error {
	return s.repo.Delete(id)
}
