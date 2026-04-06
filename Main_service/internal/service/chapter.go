package service

import (
	"errors"

	"github.com/baigel/lms/main-service/internal/model"
	"github.com/baigel/lms/main-service/internal/repository"
	"github.com/baigel/lms/main-service/pkg/apperror"
	"gorm.io/gorm"
)

type ChapterService interface {
	CreateChapter(chapter *model.Chapter) error
	GetChapterByID(id uint) (*model.Chapter, error)
	GetChaptersByCourseID(courseID uint) ([]model.Chapter, error)
	UpdateChapter(chapter *model.Chapter) (*model.Chapter, error)
	DeleteChapter(id uint) error
}

type chapterService struct {
	repo repository.ChapterRepository
}

func NewChapterService(repo repository.ChapterRepository) ChapterService {
	return &chapterService{repo: repo}
}

func (s *chapterService) CreateChapter(chapter *model.Chapter) error {
	if err := s.repo.Create(chapter); err != nil {
		return apperror.Internal("failed to create chapter", err)
	}
	return nil
}

func (s *chapterService) GetChapterByID(id uint) (*model.Chapter, error) {
	chapter, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFound("chapter not found")
		}
		return nil, apperror.Internal("failed to get chapter", err)
	}
	return chapter, nil
}

func (s *chapterService) GetChaptersByCourseID(courseID uint) ([]model.Chapter, error) {
	chapters, err := s.repo.GetAllByCourseID(courseID)
	if err != nil {
		return nil, apperror.Internal("failed to get chapters", err)
	}
	return chapters, nil
}

func (s *chapterService) UpdateChapter(chapter *model.Chapter) (*model.Chapter, error) {
	if err := s.repo.Update(chapter); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFound("chapter not found")
		}
		return nil, apperror.Internal("failed to update chapter", err)
	}

	updated, err := s.repo.GetByID(chapter.ID)
	if err != nil {
		return nil, apperror.Internal("failed to fetch updated chapter", err)
	}
	return updated, nil
}

func (s *chapterService) DeleteChapter(id uint) error {
	if err := s.repo.Delete(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.NotFound("chapter not found")
		}
		return apperror.Internal("failed to delete chapter", err)
	}
	return nil
}
