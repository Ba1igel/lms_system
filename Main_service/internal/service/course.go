package service

import (
	"github.com/baigel/lms/main-service/internal/model"
	"github.com/baigel/lms/main-service/internal/repository"
)

type CourseService interface {
	CreateCourse(course *model.Course) error
	GetCourseByID(id uint) (*model.Course, error)
	GetAllCourses() ([]model.Course, error)
	UpdateCourse(course *model.Course) error
	DeleteCourse(id uint) error
}

type courseService struct {
	repo repository.CourseRepository
}

func NewCourseService(repo repository.CourseRepository) CourseService {
	return &courseService{repo: repo}
}

func (s *courseService) CreateCourse(course *model.Course) error {
	return s.repo.Create(course)
}

func (s *courseService) GetCourseByID(id uint) (*model.Course, error) {
	return s.repo.GetByID(id)
}

func (s *courseService) GetAllCourses() ([]model.Course, error) {
	return s.repo.GetAll()
}

func (s *courseService) UpdateCourse(course *model.Course) error {
	return s.repo.Update(course)
}

func (s *courseService) DeleteCourse(id uint) error {
	return s.repo.Delete(id)
}
