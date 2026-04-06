package service

import (
	"github.com/baigel/lms/main-service/internal/model"
	"github.com/baigel/lms/main-service/internal/repository"
)

type ChapterService interface {
	CreateChapter(chapter *model.Chapter) error
	GetChapterByID(id uint) (*model.Chapter, error)
	GetChaptersByCourseID(courseID uint) ([]model.Chapter, error)
	UpdateChapter(chapter *model.Chapter) error
	DeleteChapter(id uint) error
}

type chapterService struct {
	repo repository.ChapterRepository
}

func NewChapterService(repo repository.ChapterRepository) ChapterService {
	return &chapterService{repo: repo}
}

func (s *chapterService) CreateChapter(chapter *model.Chapter) error {
	return s.repo.Create(chapter)
}

func (s *chapterService) GetChapterByID(id uint) (*model.Chapter, error) {
	return s.repo.GetByID(id)
}

func (s *chapterService) GetChaptersByCourseID(courseID uint) ([]model.Chapter, error) {
	return s.repo.GetAllByCourseID(courseID)
}

func (s *chapterService) UpdateChapter(chapter *model.Chapter) error {
	return s.repo.Update(chapter)
}

func (s *chapterService) DeleteChapter(id uint) error {
	return s.repo.Delete(id)
}
