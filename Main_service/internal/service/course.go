package service

import (
	"errors"

	"github.com/baigel/lms/main-service/internal/model"
	"github.com/baigel/lms/main-service/internal/repository"
	"github.com/baigel/lms/main-service/pkg/apperror"
	"github.com/baigel/lms/main-service/pkg/pagination"
	"gorm.io/gorm"
)

type CourseService interface {
	CreateCourse(course *model.Course) error
	GetCourseByID(id uint) (*model.Course, error)
	GetAllCourses(p pagination.Params) ([]model.Course, int64, error)
	UpdateCourse(id uint, updates map[string]any) (*model.Course, error)
	DeleteCourse(id uint) error
}

type courseService struct {
	repo repository.CourseRepository
}

func NewCourseService(repo repository.CourseRepository) CourseService {
	return &courseService{repo: repo}
}

func (s *courseService) CreateCourse(course *model.Course) error {
	if err := s.repo.Create(course); err != nil {
		return apperror.Internal("failed to create course", err)
	}
	return nil
}

func (s *courseService) GetCourseByID(id uint) (*model.Course, error) {
	course, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFound("course not found")
		}
		return nil, apperror.Internal("failed to get course", err)
	}
	return course, nil
}

func (s *courseService) GetAllCourses(p pagination.Params) ([]model.Course, int64, error) {
	courses, total, err := s.repo.GetAll(p)
	if err != nil {
		return nil, 0, apperror.Internal("failed to get courses", err)
	}
	return courses, total, nil
}

func (s *courseService) UpdateCourse(id uint, updates map[string]any) (*model.Course, error) {
	if _, err := s.GetCourseByID(id); err != nil {
		return nil, err
	}

	if err := s.repo.Update(id, updates); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFound("course not found")
		}
		return nil, apperror.Internal("failed to update course", err)
	}

	updated, err := s.repo.GetByID(id)
	if err != nil {
		return nil, apperror.Internal("failed to fetch updated course", err)
	}
	return updated, nil
}

func (s *courseService) DeleteCourse(id uint) error {
	if err := s.repo.Delete(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.NotFound("course not found")
		}
		return apperror.Internal("failed to delete course", err)
	}
	return nil
}
